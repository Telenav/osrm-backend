# osrm-api-oasisd docker

Image within build `osrm-api-oasisd` binaries, excluding generating pre-processing data by `place-connectivity-gen`, or start web service of `oasisd`.  

## Build Image

- [Dockerfile](./Dockerfile)

```bash
$ cd docker-orchestration/osrm-api-oasisd
$ docker build -t osrm-api-oasisd .   
```

## Build data

`docker run [EXTRA COMMAND] telenavmap/osrm-api-oasisd build_oasis_from_json [JSON_FILE]`  

### Parameters
- Mandatory parameters  
  - JSON_FILE: could be http url or local file  
- Optional commands 
  - `--mount "type=bind,src=$(pwd),dst=/workspace/mnt"`: mount for local PBF file input or built data output  
  - `-osrm OSRM_ENDPOINT`: osrm endpoint to calculate shortest path, if not assigned then with use great circle distance

### Samples
```bash
$ pwd
/github.com/Telenav/osrm-backend/docker-orchestration/osrm-api-oasisd

$ ls -lk
-rw-r--r--  1 ngxuser  staff       163 Jul 20 15:21 Dockerfile
-rw-r--r--  1 ngxuser  staff       319 Jul 20 15:21 README.md
-rwxr-xr-x  1 ngxuser  staff      1583 Jul 20 15:21 docker-entrypoint.sh
-rw-------@ 1 ngxuser  staff  18403706 Mar 30 09:15 us.json

$ docker run --mount "type=bind,src=$(pwd),dst=/workspace/mnt" osrm-api-oasisd:latest build_oasis_from_json -i us.json

$ ls -lk oasisdata
total 6812660
-rw-r--r--  1 ngxuser  staff     3532902 Jul 20 15:40 cellID2PointIDs.gob
-rw-r--r--  1 ngxuser  staff         428 Jul 20 15:44 connectivity_map_statistic.json
-rw-r--r--  1 ngxuser  staff  6956466286 Jul 20 15:44 id2nearbyidsmap.gob
-rw-r--r--  1 ngxuser  staff      639266 Jul 20 15:40 pointID2Location.gob
```