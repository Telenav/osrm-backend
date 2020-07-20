#!/bin/bash -xe
BUILD_PATH=${BUILD_PATH:="/osrm-build"}
DATA_PATH=${DATA_PATH:="/osrm-data"}
LOG_PATH=${LOG_PATH:="/osrm-logs"}
INPUT_STATION_DATA_JSON_FORMAT=input.json
MOUNT_PATH=/workspace/mnt # mounted host folder
OASIS_DATA_RELATIVE_PATH="oasisdata"

if [ "$1" = 'build_oasis_from_json' ]; then

  # parse input parameters, should align with integration/cmd/place-connectivity-gen/flags.go
  JSON_FILE=
  OSRM_ENDPOINT=

  # ship first parameter
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
