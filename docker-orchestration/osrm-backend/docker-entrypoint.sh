#!/bin/bash -xe
BUILD_PATH=${BUILD_PATH:="/osrm-build"}
DATA_PATH=${DATA_PATH:="/osrm-data"}
OSRM_EXTRA_COMMAND="-l DEBUG"
OSRM_ROUTED_STARTUP_COMMAND=" -a MLD --max-table-size 8000 "
MAPDATA_NAME_WITH_SUFFIX=map
PBF_FILE_SUFFIX=".osm.pbf"
WAYID2NODEIDS_MAPPING_FILE=wayid2nodeids.csv
WAYID2NODEIDS_MAPPING_FILE_COMPRESSED=${WAYID2NODEIDS_MAPPING_FILE}.snappy

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
  ${BUILD_PATH}/osrm-traffic-updater -c ${TRAFFIC_PROXY_IP} -m ${WAYID2NODEIDS_MAPPING_FILE_COMPRESSED} -f ${TRAFFIC_FILE} -map ${MAP_PROVIDER} -traffic ${TRAFFIC_PROVIDER} -region ${REGION} ${BLOCKING_ONLY}
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
  ${BUILD_PATH}/osrm-extract $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX} -p ${BUILD_PATH}/${PROFILE_LUA} -d ${DATA_VERSION} ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-partition $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-customize $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  if [ "${PBF_SOURCE}" = "unidb" ]; then
    ${BUILD_PATH}/wayid2nodeid-extractor -i $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX} -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -b=true
    ${BUILD_PATH}/snappy -i $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE_COMPRESSED}
  fi
  ls -lh ${DATA_PATH}/

  # clean source pbf and temp files
  rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}${PBF_FILE_SUFFIX}
  rm -f $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE}
  if [ "${KEEP_TEMP_OSRM_FILES}" != "true" ]; then # set KEEP_TEMP_OSRM_FILES explicitly by env vars.
    rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm
  fi
  ls -lh ${DATA_PATH}/

  # export compiled mapdata to mounted path for publishing 
  COMPILED_DATA_EXPORT_PATH=/compiled-data
  mv ${DATA_PATH}/* ${COMPILED_DATA_EXPORT_PATH}/
  chmod 777 ${COMPILED_DATA_EXPORT_PATH}/*

else
  exec "$@"
fi
