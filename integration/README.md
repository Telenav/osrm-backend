# OSRM Integration
All OSRM integration related codes will be organized together in this folder, e.g. traffic integration, command line tools, etc. Mostly they're `Go` codes, require at least `go1.13`. 

## osrm-traffic-updater
Command line tool [cmd/osrm-traffic-updater](cmd/osrm-traffic-updater/) is designed for pull traffic data from **traffic-proxy(Telenav)** then dump to OSRM required `traffic.csv` for `osrm-customize`. Refer to [OSRM with Telenav Traffic Design](doc/osrm-with-telenav-traffic.md) and [OSRM Traffic](https://github.com/Project-OSRM/osrm-backend/wiki/Traffic) for more details.        

- RPC Protocol
[proxy.proto](proxy.proto)


## wayid2nodeids_extractor
Command line tool for extract wayid to nodeids mapping from PBF. Code in [cmd/wayid2nodeid-extractor](cmd/wayid2nodeid-extractor/).        

## nodes2way-builder
Command line tool for build nodes2way mapping DB. Code in [cmd/nodes2way-builder](cmd/nodes2way-builder/).        

## nodes2way-cli
Command line tool for querying wayids from nodeids in DB. Code in [cmd/nodes2way-cli](cmd/nodes2way-cli/).        

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
I0121 15:01:55.282673    7124 contents.go:60] Loaded from nevada-latest.osrm
I0121 15:01:55.285059    7124 contents.go:61]   OSRN v5.22.0
I0121 15:01:55.285161    7124 contents.go:63]   nodes meta 1092695 count 1092695
I0121 15:01:55.285172    7124 contents.go:65]     node[0] {-120011751 39443340 26798725}
I0121 15:01:55.285192    7124 contents.go:65]     node[1] {-120017543 39440794 26798726}
I0121 15:01:55.285202    7124 contents.go:65]     node[2] {-120031913 39431791 26798727}
I0121 15:01:55.285228    7124 contents.go:65]     node[3] {-120035895 39423872 26798728}
I0121 15:01:55.285256    7124 contents.go:65]     node[4] {-120030001 39416698 26798729}
I0121 15:01:55.285266    7124 contents.go:68]   barriers meta 85 count 85
I0121 15:01:55.285292    7124 contents.go:70]     barrier[0] 68738
I0121 15:01:55.285321    7124 contents.go:70]     barrier[1] 127356
I0121 15:01:55.285338    7124 contents.go:70]     barrier[2] 205601
I0121 15:01:55.285348    7124 contents.go:70]     barrier[3] 331924
I0121 15:01:55.285410    7124 contents.go:70]     barrier[4] 351581
I0121 15:01:55.285439    7124 contents.go:73]   traffic_lights meta 3821 count 3821
I0121 15:01:55.285466    7124 contents.go:75]     traffic_lights[0] 6584
I0121 15:01:55.285493    7124 contents.go:75]     traffic_lights[1] 6596
I0121 15:01:55.285524    7124 contents.go:75]     traffic_lights[2] 6600
I0121 15:01:55.285548    7124 contents.go:75]     traffic_lights[3] 15405
I0121 15:01:55.285576    7124 contents.go:75]     traffic_lights[4] 15439
I0121 15:01:55.285600    7124 contents.go:78]   edges meta 1179277 count 1179277
I0121 15:01:55.285650    7124 contents.go:80]     edges[0] {0 503728 24 24 54.78657 {2147483647 false} 44793 {true false false false false true false {true false false 0 2} 0 0}}
I0121 15:01:55.285726    7124 contents.go:80]     edges[1] {503729 0 21 21 49.09447 {2147483647 false} 44793 {true false false false false true false {true false false 0 2} 0 0}}
I0121 15:01:55.285753    7124 contents.go:80]     edges[2] {503726 1 50 50 116.59047 {2147483647 false} 44793 {true false false false false true false {true false false 0 2} 0 0}}
I0121 15:01:55.285785    7124 contents.go:80]     edges[3] {1 618170 23 23 53.611294 {2147483647 false} 44793 {true false false false false true false {true false false 0 2} 0 0}}
I0121 15:01:55.285814    7124 contents.go:80]     edges[4] {503714 2 54 54 126.03433 {2147483647 false} 44793 {true false false false false true false {true false false 0 2} 0 0}}
I0121 15:01:55.285840    7124 contents.go:83]   annotations meta 158342 count 158342
I0121 15:01:55.285851    7124 contents.go:85]     annotations[0] {5 65535 1 1 false}
I0121 15:01:55.285899    7124 contents.go:85]     annotations[1] {0 65535 0 1 false}
I0121 15:01:55.285924    7124 contents.go:85]     annotations[2] {10 65535 0 1 false}
I0121 15:01:55.285951    7124 contents.go:85]     annotations[3] {0 65535 0 1 false}
I0121 15:01:55.286009    7124 contents.go:85]     annotations[4] {15 65535 0 1 false}
I0121 15:01:55.286492    7124 contents.go:52] Loaded from nevada-latest.osrm.timestamp
I0121 15:01:55.286525    7124 contents.go:53]   OSRN v5.22.0
I0121 15:01:55.286539    7124 contents.go:55]   timestamp(a.k.a. data_version) meta 20 count 20
I0121 15:01:55.286549    7124 contents.go:57]   timestamp(a.k.a. data_version) 2019-01-24T21:15:02Z
```
