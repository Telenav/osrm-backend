# OSRM Traffic Updater
The **OSRM Traffic Updater** is designed for pull traffic data from **Traffic Proxy(Telenav)** then dump to OSRM required `traffic.csv`. Refer to [OSRM with Telenav Traffic Design](../docs/design/osrm-with-telenav-traffic.md) and [OSRM Traffic](https://github.com/Project-OSRM/osrm-backend/wiki/Traffic) for more details.        

## RPC Protocol
See [proxy.proto](proxy.proto) for details.    

## Requirements
- `go1.13`

## Tools

### osrm-traffic-updater
Pull telenav traffic from `traffic-proxy` and generate OSRM required `traffic.csv` for `osrm-customize`.    

### wayid2nodeids_extractor
Extract wayid to nodeids mapping from PBF.

### snappy
Command line tool for [snappy](github.com/golang/snappy) compression.    
