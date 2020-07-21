#!/bin/bash -xe
BUILD_PATH=${BUILD_PATH:="/osrm-build"}
DATA_PATH=${DATA_PATH:="/osrm-data"}
OSRM_EXTRA_COMMAND="-l DEBUG"
OSRM_ROUTED_STARTUP_COMMAND=" -a MLD --max-table-size 8000 "
MAPDATA_NAME_WITH_SUFFIX=map
PBF_FILE_SUFFIX=".osm.pbf"
SNAPPY_SUFFIX=".snappy"
WAYID2NODEIDS_MAPPING_FILE=wayid2nodeids.csv
NODES2WAY_DB_FILE="nodes2way.db"
INPUT_STATION_DATA_JSON_FORMAT=input.json
MOUNT_PATH=/workspace/mnt # mounted host folder
OASIS_DATA_RELATIVE_PATH="oasisdata"

_sig() {
  kill -TERM $child 2>/dev/null
}

if [ "$1" = 'routed_startup' ] || [ "$1" = 'routed_blocking_traffic_startup' ]; then
  #trap _sig SIGKILL SIGTERM SIGHUP SIGINT EXIT

  TRAFFIC_FILE=traffic.csv
  TRAFFIC_PROXY_IP=${2:-"10.189.102.81"}
  REGION=${3}
  MAP_PROVIDER=${4}
  TRAFFIC_PROVIDER=${5}
  if [ "$1" = 'routed_blocking_traffic_startup' ]; then
    BLOCKING_ONLY="-blocking-only"
  fi
  if [ "$6" = 'incremental' ]; then
    INCREMENTAL_CUSTOMIZE="--incremental true"
  fi

  cd ${DATA_PATH}
  ${BUILD_PATH}/osrm-traffic-updater -c ${TRAFFIC_PROXY_IP} -m ${WAYID2NODEIDS_MAPPING_FILE}${SNAPPY_SUFFIX} -f ${TRAFFIC_FILE} -map ${MAP_PROVIDER} -traffic ${TRAFFIC_PROVIDER} -region ${REGION} ${BLOCKING_ONLY}
  ls -lh
  ${BUILD_PATH}/osrm-customize ${MAPDATA_NAME_WITH_SUFFIX}.osrm  --segment-speed-file ${TRAFFIC_FILE} ${OSRM_EXTRA_COMMAND} ${INCREMENTAL_CUSTOMIZE}
  ${BUILD_PATH}/osrm-routed ${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_ROUTED_STARTUP_COMMAND} &
  child=$!
  wait "$child"

elif [ "$1" = 'routed_no_traffic_startup' ]; then
  #trap _sig SIGKILL SIGTERM SIGHUP SIGINT EXIT

  cd ${DATA_PATH}
  ${BUILD_PATH}/osrm-routed ${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_ROUTED_STARTUP_COMMAND} &
  child=$!
  wait "$child"

elif [ "$1" = 'compile_mapdata' ]; then
  #trap _sig SIGKILL SIGTERM SIGHUP SIGINT EXIT

  PROFILE_LUA=${2} # e.g., "profiles/car.lua"
  PBF_FILE_URL=${3}
  PBF_SOURCE=${4:-"osm"}
  if [ x"${DATA_VERSION}" = x ]; then  # set DATA_VERSION explicitly by env vars. Use the PBF_FILE_URL if not set.
    DATA_VERSION=${PBF_FILE_URL}
  fi

  curl -sSL -f ${PBF_FILE_URL} > $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX}

  # osrm extract/partition/customize
  ${BUILD_PATH}/osrm-extract $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX} -p ${BUILD_PATH}/${PROFILE_LUA} -d ${DATA_VERSION} ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-partition $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-customize $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  
  # extract way,node,node,... from PBF
  ${BUILD_PATH}/wayid2nodeid-extractor -i $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX} -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -mapsource ${PBF_SOURCE}
  ${BUILD_PATH}/snappy -i $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE}${SNAPPY_SUFFIX}
  
  # build fromNode,toNode->way DB
  NODES2WAY_DB_FILE_BUILDING_PATH=${DATA_PATH}/${NODES2WAY_DB_FILE}
  if [ x"${SHM_PATH}" != x ]; then # set SHM explicitly by env vars, then use it to speed up the DB building process
    NODES2WAY_DB_FILE_BUILDING_PATH=${SHM_PATH}/${NODES2WAY_DB_FILE}
  fi
  ${BUILD_PATH}/nodes2way-builder -alsologtostderr -v 2 -snappy-compressed=false -i $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -o ${NODES2WAY_DB_FILE_BUILDING_PATH}
  ${BUILD_PATH}/snappy -i ${NODES2WAY_DB_FILE_BUILDING_PATH} -o ${DATA_PATH}/${NODES2WAY_DB_FILE}${SNAPPY_SUFFIX}
  ls -lh ${DATA_PATH}/

  # clean source pbf and temp files
  rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX}
  rm -f $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE}
  rm -f ${NODES2WAY_DB_FILE_BUILDING_PATH}
  if [ "${KEEP_TEMP_OSRM_FILES}" != "true" ]; then # set KEEP_TEMP_OSRM_FILES explicitly by env vars.
    rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm
  fi
  ls -lh ${DATA_PATH}/

  # export compiled mapdata to mounted path for publishing 
  COMPILED_DATA_EXPORT_PATH=/compiled-data
  mv ${DATA_PATH}/* ${COMPILED_DATA_EXPORT_PATH}/
  chmod 777 ${COMPILED_DATA_EXPORT_PATH}/*

elif [ "$1" = 'build_oasis_from_json' ]; then

  # parse input parameters, should align with integration/cmd/place-connectivity-gen/flags.go
  JSON_FILE=
  OSRM_ENDPOINT=

  # skip first parameter
  shift 1
  while [[ "$#" -gt 0 ]]; do
    case $1 in
      -i|--input) JSON_FILE="$2"; shift ;;
      -osrm) OSRM_ENDPOINT="$2"; shift ;;
      *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
  done

  # create output folder
  mkdir -p ${MOUNT_PATH}/${OASIS_DATA_RELATIVE_PATH}

  # construct path for input
  if [[ ${JSON_FILE} = http* ]]; then 
    curl -sSL -f ${JSON_FILE} > ${DATA_PATH}/${INPUT_STATION_DATA_JSON_FORMAT}
  else 
    cp ${MOUNT_PATH}/${JSON_FILE} ${DATA_PATH}/${INPUT_STATION_DATA_JSON_FORMAT}
  fi

  # generate parameters for oasis preprocessing
  OASIS_PREPROCESSING_PARAMETERS=
  if [ -z "$osrm" ]
  then
    OASIS_PREPROCESSING_PARAMETERS="${OASIS_PREPROCESSING_PARAMETERS} -i ${DATA_PATH}/${INPUT_STATION_DATA_JSON_FORMAT} -o ${MOUNT_PATH}/${OASIS_DATA_RELATIVE_PATH}"
  else
    OASIS_PREPROCESSING_PARAMETERS="${OASIS_PREPROCESSING_PARAMETERS} -i ${DATA_PATH}/${INPUT_STATION_DATA_JSON_FORMAT} -o ${MOUNT_PATH}/${OASIS_DATA_RELATIVE_PATH} -osrm ${OSRM_ENDPOINT}"
  fi

  # generate data
  ${BUILD_PATH}/place-connectivity-gen ${OASIS_PREPROCESSING_PARAMETERS}

else
  exec "$@"
fi
