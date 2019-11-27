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

