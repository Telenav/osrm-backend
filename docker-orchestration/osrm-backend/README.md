# telenav osrm-backend docker
Image within built osrm binaries(`osrm-extract/osrm-partition/osrm-customize/...`) and running dependencies. It can be used to **compile data** or **startup routed**.      

## Image

```bash
$ docker push telenavmap/osrm-backend 
```

## Run 
### Run with osrm-data outside

```bash
# prepare compiled data first (sample data: california)
$ cd osrm-data
$ ll -lh
total 243M
-rw-r--r-- 1 mapuser appuser  21M Jun 17 23:37 map.osrm.cell_metrics
-rw-r--r-- 1 mapuser appuser 193K Jun 17 23:34 map.osrm.cells
-rw-r--r-- 1 mapuser appuser 1.9M Jun 17 23:34 map.osrm.cnbg
-rw-r--r-- 1 mapuser appuser 1.9M Jun 17 23:34 map.osrm.cnbg_to_ebg
-rw-r--r-- 1 mapuser appuser  68K Jun 17 23:37 map.osrm.datasource_names
-rw-r--r-- 1 mapuser appuser 9.8M Jun 17 23:34 map.osrm.ebg
-rw-r--r-- 1 mapuser appuser 2.8M Jun 17 23:34 map.osrm.ebg_nodes
-rw-r--r-- 1 mapuser appuser 2.9M Jun 17 23:34 map.osrm.edges
-rw-r--r-- 1 mapuser appuser 2.7M Jun 17 23:34 map.osrm.enw
-rwx------ 1 mapuser appuser 5.6M Jun 17 23:34 map.osrm.fileIndex
-rw-r--r-- 1 mapuser appuser 7.3M Jun 17 23:37 map.osrm.geometry
-rw-r--r-- 1 mapuser appuser 1.1M Jun 17 23:34 map.osrm.icd
-rw-r--r-- 1 mapuser appuser 5.0K Jun 17 23:34 map.osrm.maneuver_overrides
-rw-r--r-- 1 mapuser appuser  11M Jun 17 23:37 map.osrm.mldgr
-rw-r--r-- 1 mapuser appuser 218K Jun 17 23:34 map.osrm.names
-rw-r--r-- 1 mapuser appuser 4.0M Jun 17 23:34 map.osrm.nbg_nodes
-rw-r--r-- 1 mapuser appuser 1.8M Jun 17 23:34 map.osrm.partition
-rw-r--r-- 1 mapuser appuser 6.0K Jun 17 23:34 map.osrm.properties
-rw-r--r-- 1 mapuser appuser  29K Jun 17 23:34 map.osrm.ramIndex
-rw-r--r-- 1 mapuser appuser 4.0K Jun 17 23:34 map.osrm.restrictions
-rw-r--r-- 1 mapuser appuser 3.5K Jun 17 23:34 map.osrm.timestamp
-rw-r--r-- 1 mapuser appuser 5.5K Jun 17 23:34 map.osrm.tld
-rw-r--r-- 1 mapuser appuser 8.0K Jun 17 23:34 map.osrm.tls
-rw-r--r-- 1 mapuser appuser 836K Jun 17 23:34 map.osrm.turn_duration_penalties
-rw-r--r-- 1 mapuser appuser 4.9M Jun 17 23:34 map.osrm.turn_penalties_index
-rw-r--r-- 1 mapuser appuser 836K Jun 17 23:34 map.osrm.turn_weight_penalties
$ cd ..

# pull & run
$ docker pull telenavmap/osrm-backend:no.63-20191112-master-telenav-845f953
$ docker run -d -p 5000:5000  "src=$(pwd)/osrm-data,dst=/osrm-data,type=bind" telenavmap/osrm-backend:no.63-20191112-master-telenav-845f953 routed_no_traffic_startup 
```

## Example By Manual
- [Build Berlin Server with OSM data](./example-berlin-osm.md)

