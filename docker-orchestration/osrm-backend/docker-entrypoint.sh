#!/bin/bash -xe
BUILD_PATH=${BUILD_PATH:="/osrm-build"}
DATA_PATH=${DATA_PATH:="/osrm-data"}
OSRM_EXTRA_COMMAND="-l DEBUG"
OSRM_ROUTED_STARTUP_COMMAND=" -a MLD --max-table-size 8000 "
MAPDATA_NAME_WITH_SUFFIX=map
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

  PBF_FILE_URL=${2}
  IS_TELENAV_PBF=${3:-"false"}
  DATA_VERSION=${4:-"unset"}

  curl -sSL -f ${PBF_FILE_URL} > $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osm.pbf
  ${BUILD_PATH}/osrm-extract $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osm.pbf -p ${BUILD_PATH}/profiles/car.lua -d ${DATA_VERSION} ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-partition $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/osrm-customize $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm ${OSRM_EXTRA_COMMAND}
  ${BUILD_PATH}/wayid2nodeid-extractor -i $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osm.pbf -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE}  -b=${IS_TELENAV_PBF}
  ${BUILD_PATH}/snappy -i $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE} -o $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE_COMPRESSED}
  ls -lh $DATA_PATH/

  # clean source pbf and temp files
  rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osm.pbf
  rm -f $DATA_PATH/${MAPDATA_NAME_WITH_SUFFIX}.osrm
  rm -f $DATA_PATH/${WAYID2NODEIDS_MAPPING_FILE}

  # export compiled mapdata to mounted path for publishing 
  SAVE_DATA_PACKAGE_PATH=/save-data
  mv ${DATA_PATH}/* ${SAVE_DATA_PACKAGE_PATH}/
  chmod 777 ${SAVE_DATA_PACKAGE_PATH}/*

else
  exec "$@"
fi
