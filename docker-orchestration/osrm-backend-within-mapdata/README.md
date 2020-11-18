# osrm-backend-within-mapdata
Image based on [osrm-backend docker image](../osrm-backend/) and put compiled mapdata inside.          
NOTE: It's a temporary workaround for easily run in k8s. It's NOT a good idea to put mapdata in image directly since the map data is too big. For long-term, discussing in [#93](https://github.com/Telenav/osrm-backend/issues/93).      

## Build Image
- [Dockerfile](./Dockerfile)

```bash
$ cd docker-orchestration/osrm-backend-within-mapdata/
$ 
$ # TODO: copy OSRM built files(*.osrm.*) to subfolder `map` of this folder
$
$ ls -lh map/
ll -lh
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
$
$ docker build -t osrm-backend-within-mapdata --build-arg FROM_TAG=${YOUR_OSRM_BACKEND_IMAGE_TAG}  .   
```

## Pull Image
```bash
$ docker pull telenavmap/osrm-backend-within-mapdata
```

## Run osrm-routed

### No traffic
```bash
$ docker run -d -p 5000:5000 telenavmap/osrm-backend-within-mapdata routed_no_traffic_startup
```

### With traffic

```bash
$ docker run -d -p 5000:5000 [ENV_OPTIONS] telenavmap/osrm-backend-within-mapdata routed_startup
```

| Option | Type | Mandatory/Optional | Description |
|--------|------|--------------------|-------------|
|`routed_startup` | Param | Mandatory | `routed_startup`: fetch and customize full region traffic|
| `ENABLE_INCREMENTAL_CUSTOMIZE` | ENV | Optional | only customize cells that traffic touches to reduce time, defaultly customize all cells. Set `ENABLE_INCREMENTAL_CUSTOMIZE=true` by env to enable it. |
| `CUSTOMIZE_THREADS` | ENV | Optional | Maximum threads to run `osrm-customize`. |  
| `TRAFFIC_PROXY_ENDPOINT` | ENV | Optional | `traffic-proxy` endpoint to connect, `10.189.102.74:10086` by default. |
| `MAP_PROVIDER` | ENV | Optional | map provider, `osm` by default. | 
| `TRAFFIC_PROVIDER` | ENV | Optional | traffic provider. | 
| `REGION` | ENV | Optional | which region of traffic to fetch, `na` by default. |
| `FETCH_TRAFFIC_EXTRA_ARGS` | ENV | Optional | Additional args for `osrm-traffic-updater`, for debugging purpose, e.g., `-v 2` | 

```bash
$ # example
$ docker run -d -p 5000:5000 -e TRAFFIC_PROXY_ENDPOINT=${YOUR_TRAFFIC_PROXY_ENDPOINT} -e TRAFFIC_PROVIDER=${YOUR_TRAFFIC_PROVIDER} -e REGION=${YOUR_REGION} telenavmap/osrm-backend-within-mapdata routed_startup
```

## Run osrm-rankd


```bash
$ docker run -d -p 5001:5000 --shm-size=64g telenavmap/osrm-backend-within-mapdata rankd_startup -osrm 127.0.0.1:5000
$ 
$ # get route with ways in annotations
$ curl "http://127.0.0.1:5001/route/v1/driving/-115.130556,36.126124;-115.151499,-36.163968?annotations=true"
```