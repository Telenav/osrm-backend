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
