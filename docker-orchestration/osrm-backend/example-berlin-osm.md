# Build OSRM server based on Berlin OSM data

- Generate docker image
```bash
$ DOCKER_BUILDKIT=1 docker build --no-cache -t telenavmap/osrm-backend --build-arg GIT_COMMIT=master --build-arg CMAKE_BUILD_EXTRA_ARGS="-- -j" .
```

- Generate OSRM data
```bash
$ mkdir -p compiled-data
$
$ docker run -it --mount "src=$(pwd)/compiled-data,dst=/compiled-data,type=bind" telenavmap/osrm-backend compile_mapdata "profiles/car.lua" "https://download.geofabrik.de/europe/germany/berlin-latest.osm.pbf" 
$ 
$ # optional keep the temporarily .osrm file by '-e KEEP_TEMP_OSRM_FILES=true'
$ # optional set data_version explicitly by '-e DATA_VERSION=YOUR_DATA_VERSION '
$ # e.g., docker run -it --mount "src=$(pwd)/compiled-data,dst=/compiled-data,type=bind" -e KEEP_TEMP_OSRM_FILES=true -e DATA_VERSION=YOUR_DATA_VERSION telenavmap/osrm-backend compile_mapdata "profiles/car.lua" "https://download.geofabrik.de/europe/germany/berlin-latest.osm.pbf" 
$ 
$ ls -lh compiled-data/
total 88M
-rwxrwxrwx 1 root root  23M Mar 16 20:03 map.osrm.cell_metrics
-rwxrwxrwx 1 root root 214K Mar 16 20:03 map.osrm.cells
-rwxrwxrwx 1 root root 2.1M Mar 16 20:03 map.osrm.cnbg
-rwxrwxrwx 1 root root 2.1M Mar 16 20:03 map.osrm.cnbg_to_ebg
-rwxrwxrwx 1 root root  68K Mar 16 20:03 map.osrm.datasource_names
-rwxrwxrwx 1 root root  11M Mar 16 20:03 map.osrm.ebg
-rwxrwxrwx 1 root root 3.0M Mar 16 20:03 map.osrm.ebg_nodes
-rwxrwxrwx 1 root root 3.2M Mar 16 20:03 map.osrm.edges
-rwxrwxrwx 1 root root 2.9M Mar 16 20:03 map.osrm.enw
-rwxrwxrwx 1 root root 6.1M Mar 16 20:03 map.osrm.fileIndex
-rwxrwxrwx 1 root root 8.0M Mar 16 20:03 map.osrm.geometry
-rwxrwxrwx 1 root root 1.2M Mar 16 20:03 map.osrm.icd
-rwxrwxrwx 1 root root 5.0K Mar 16 20:03 map.osrm.maneuver_overrides
-rwxrwxrwx 1 root root  12M Mar 16 20:03 map.osrm.mldgr
-rwxrwxrwx 1 root root 218K Mar 16 20:03 map.osrm.names
-rwxrwxrwx 1 root root 4.4M Mar 16 20:03 map.osrm.nbg_nodes
-rwxrwxrwx 1 root root 2.0M Mar 16 20:03 map.osrm.partition
-rwxrwxrwx 1 root root 6.0K Mar 16 20:03 map.osrm.properties
-rwxrwxrwx 1 root root  31K Mar 16 20:03 map.osrm.ramIndex
-rwxrwxrwx 1 root root 4.0K Mar 16 20:03 map.osrm.restrictions
-rwxrwxrwx 1 root root 4.0K Mar 16 20:03 map.osrm.timestamp
-rwxrwxrwx 1 root root 5.5K Mar 16 20:03 map.osrm.tld
-rwxrwxrwx 1 root root 8.0K Mar 16 20:03 map.osrm.tls
-rwxrwxrwx 1 root root 922K Mar 16 20:03 map.osrm.turn_duration_penalties
-rwxrwxrwx 1 root root 5.4M Mar 16 20:03 map.osrm.turn_penalties_index
-rwxrwxrwx 1 root root 922K Mar 16 20:03 map.osrm.turn_weight_penalties
```

- Start OSRM server
```bash
$ docker run -d -p 5000:5000 --mount "src=$(pwd)/compiled-data,dst=/osrm-data,type=bind" telenavmap/osrm-backend routed_no_traffic_startup
d26598e503e8dc57088b131559c3c94279d08c16a066147abcdf1c5d088d7704
$
$ docker logs -f d26598e503e8dc57088b131559c3c94279d08c16a066147abcdf1c5d088d7704
+ BUILD_PATH=/osrm-build
+ DATA_PATH=/osrm-data
+ OSRM_EXTRA_COMMAND='-l DEBUG'
+ OSRM_ROUTED_STARTUP_COMMAND=' -a MLD --max-table-size 8000 '
+ MAPDATA_NAME_WITH_SUFFIX=map
+ PBF_FILE_SUFFIX=.osm.pbf
+ WAYID2NODEIDS_MAPPING_FILE=wayid2nodeids.csv
+ WAYID2NODEIDS_MAPPING_FILE_COMPRESSED=wayid2nodeids.csv.snappy
+ '[' routed_no_traffic_startup = routed_startup ']'
+ '[' routed_no_traffic_startup = routed_blocking_traffic_startup ']'
+ '[' routed_no_traffic_startup = routed_no_traffic_startup ']'
+ cd /osrm-data
+ child=7
+ wait 7
+ /osrm-build/osrm-routed map.osrm -a MLD --max-table-size 8000
[info] starting up engines, v5.22.0
[info] Threads: 8
[info] IP address: 0.0.0.0
[info] IP port: 5000
[info] http 1.1 compression handled by zlib version 1.2.8
[info] Listening on: 0.0.0.0:5000
[info] running and waiting for requests

```

- Test
```bash
curl "http://127.0.0.1:5000/route/v1/driving/13.388860,52.517037;13.397634,52.529407"
```


