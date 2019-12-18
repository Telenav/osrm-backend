# OSRM Integration
All OSRM integration related codes will be organized together in this folder, e.g. traffic integration, command line tools, etc. Mostly they're `Go` codes, require at least `go1.13`. 

## osrm-traffic-updater
Command line tool [cmd/osrm-traffic-updater](cmd/osrm-traffic-updater/) is designed for pull traffic data from **traffic-proxy(Telenav)** then dump to OSRM required `traffic.csv` for `osrm-customize`. Refer to [OSRM with Telenav Traffic Design](doc/osrm-with-telenav-traffic.md) and [OSRM Traffic](https://github.com/Project-OSRM/osrm-backend/wiki/Traffic) for more details.        

- RPC Protocol
[proxy.proto](proxy.proto)


## wayid2nodeids_extractor
Command line tool for extract wayid to nodeids mapping from PBF. Code in [cmd/wayid2nodeid-extractor](cmd/wayid2nodeid-extractor/).        

## snappy
Command line tool for [snappy](github.com/golang/snappy) compression. Code in [cmd/snappy](cmd/snappy/).  

## osrm-ranking 
Update `duration/weight` by traffic for many `alternatives`, then pick up best as result.     
- design [OSRM with Telenav Traffic Design - Alternatives Ranking](doc/osrm-with-telenav-traffic.md)     
- code [cmd/osrm-ranking](cmd/osrm-ranking)    
- monitor API: `/monitor`, e.g. `http://localhost:8080/monitor`     

## trafficproxy-cli 
Command line tool for querying traffic from `trafficproxy`. Code in [cmd/trafficproxy-cli](cmd/trafficproxy-cli/).       
Typical usage:    

```bash
# 1. get traffic for ways 
$ trafficproxy-cli -c ${TARGET_IP} -map ${MAP} -traffic ${TRAFFIC} -ways 829733412,-104489539,-129639168

# 2. get all traffic 
$ trafficproxy-cli -c ${TARGET_IP} -map ${MAP} -traffic ${TRAFFIC} -mode getall -stdout=false -dumpfile test 

# 3. steaming delta traffic 
$ trafficproxy-cli -c ${TARGET_IP} -map ${MAP} -traffic ${TRAFFIC} -mode delta -stdout=false -dumpfile test

# if want to see more running log by shell, append `-alsologtostderr` in command-line
# if want to see more running log by log files, append `-log_dir=.` in command-line
$ trafficproxy-cli -alsologtostderr ...
$ trafficproxy-cli -log_dir=. ...

# for more options, see help
$ trafficproxy-cli -h

```

## trafficcache-parallel-test
Command line tool for traffic cache test. There could be two type of traffic caches, i.e. indexed by wayID and indexed by Edge. This tool possible to run them in parallel and do some comparison test. Code in [cmd/trafficcache-parallel-test](cmd/trafficcache-parallel-test/).          

## osrm-files-extractor
Command line tool for extract and show information of OSRM generated `.osrm[.*]` files.     
Code in [cmd/osrm-files-extractor](cmd/osrm-files-extractor/).     
Typical usage:    
```bash
$ # for each `.osrm[.*]` file, we can extermine it by `tar`. E.g. 
$ tar tvf nevada-latest.osrm 
-rw-rw-r--  0 0      0           8 Jan  1  1970 osrm_fingerprint.meta
-rw-rw-r--  0 0      0           8 Jan  1  1970 /extractor/nodes.meta
-rw-rw-r--  0 0      0    18275776 Jan  1  1970 /extractor/nodes
-rw-rw-r--  0 0      0           8 Jan  1  1970 /extractor/barriers.meta
-rw-rw-r--  0 0      0         608 Jan  1  1970 /extractor/barriers
-rw-rw-r--  0 0      0           8 Jan  1  1970 /extractor/traffic_lights.meta
-rw-rw-r--  0 0      0       16144 Jan  1  1970 /extractor/traffic_lights
-rw-rw-r--  0 0      0           8 Jan  1  1970 /extractor/edges.meta
-rw-rw-r--  0 0      0    39504672 Jan  1  1970 /extractor/edges
-rw-rw-r--  0 0      0           8 Jan  1  1970 /extractor/annotations.meta
-rw-rw-r--  0 0      0     1355088 Jan  1  1970 /extractor/annotations
$ 
$ # we can use `osrm-files-extractor` to see its details.    
$ ./osrm-files-extractor -alsologtostderr -f nevada-latest.osrm -summary 5
I1218 14:12:25.868053    1747 contents.go:81] Loaded from nevada-latest.osrm
I1218 14:12:25.869828    1747 contents.go:82]   OSRN v5.22.0
I1218 14:12:25.869848    1747 contents.go:84]   nodes meta 1142236 count 1142236
I1218 14:12:25.869853    1747 contents.go:86]     node[0] {-120011751 39443340 26798725}
I1218 14:12:25.870092    1747 contents.go:86]     node[1] {-120017543 39440794 26798726}
I1218 14:12:25.870097    1747 contents.go:86]     node[2] {-120031913 39431791 26798727}
I1218 14:12:25.870100    1747 contents.go:86]     node[3] {-120035895 39423872 26798728}
I1218 14:12:25.870103    1747 contents.go:86]     node[4] {-120030001 39416698 26798729}
I1218 14:12:25.870106    1747 contents.go:89]   barriers meta 152 count 152
I1218 14:12:25.870110    1747 contents.go:91]     barrier[0] 66546
I1218 14:12:25.870114    1747 contents.go:91]     barrier[1] 196332
I1218 14:12:25.870117    1747 contents.go:91]     barrier[2] 235061
I1218 14:12:25.870122    1747 contents.go:91]     barrier[3] 316127
I1218 14:12:25.870125    1747 contents.go:91]     barrier[4] 335454
I1218 14:12:25.870128    1747 contents.go:94]   traffic_lights meta 4036 count 4036
I1218 14:12:25.870132    1747 contents.go:96]     traffic_lights[0] 6467
I1218 14:12:25.870135    1747 contents.go:96]     traffic_lights[1] 6479
I1218 14:12:25.870138    1747 contents.go:96]     traffic_lights[2] 6483
I1218 14:12:25.870140    1747 contents.go:96]     traffic_lights[3] 13948
I1218 14:12:25.870143    1747 contents.go:96]     traffic_lights[4] 13981
I1218 14:12:25.870146    1747 contents.go:99]   edges meta 1234521 count 1234521
I1218 14:12:25.870150    1747 contents.go:101]     edges[0] {0 487556 24 24 54.78657 {2147483647 false} 44494 {true false false false false true false {true false false 0 2} 0 0}}
I1218 14:12:25.870844    1747 contents.go:101]     edges[1] {487557 0 21 21 49.09447 {2147483647 false} 44494 {true false false false false true false {true false false 0 2} 0 0}}
I1218 14:12:25.870853    1747 contents.go:101]     edges[2] {487554 1 50 50 116.59047 {2147483647 false} 44494 {true false false false false true false {true false false 0 2} 0 0}}
I1218 14:12:25.870860    1747 contents.go:101]     edges[3] {1 604105 23 23 53.611294 {2147483647 false} 44494 {true false false false false true false {true false false 0 2} 0 0}}
I1218 14:12:25.870866    1747 contents.go:101]     edges[4] {487542 2 54 54 126.03433 {2147483647 false} 44494 {true false false false false true false {true false false 0 2} 0 0}}
I1218 14:12:25.870872    1747 contents.go:104]   annotations meta 169386 count 169386
I1218 14:12:25.870876    1747 contents.go:106]     annotations[0] {5 65535 1 1 false}
I1218 14:12:25.870884    1747 contents.go:106]     annotations[1] {0 65535 0 1 false}
I1218 14:12:25.870888    1747 contents.go:106]     annotations[2] {10 65535 0 1 false}
I1218 14:12:25.870892    1747 contents.go:106]     annotations[3] {15 65535 0 1 false}
I1218 14:12:25.870896    1747 contents.go:106]     annotations[4] {0 65535 0 1 false}
```
