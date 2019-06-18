# telenav osrm-backend docker
Image within built osrm binaries(`osrm-extract/osrm-partition/osrm-customize/...`) and running dependencies. It can be used to **compile data** or **startup routed**.      

## Build Image
Please use below jenkins job to build docker image within osrm binaries from source code.    

- [Jenkins Job - Build_Telenav_OSRM_Backend_Docker](https://shd-routingfp-01.telenav.cn:8443/view/OSRM/job/Build_Telenav_OSRM_Backend_Docker/)    

## Compile Data
Please use below jenkins job to compile mapdata from PBF to osrm data.       
There're two options for publish the compiled osrm data:    
- generate new docker image from current one within compile data
- export as `map.tar.gz`    

See below job for details:      
- [Jenkins Job - Compile_Mapdata_In_Telenav_OSRM_Backend_Docker](https://shd-routingfp-01.telenav.cn:8443/view/OSRM/job/Compile_Mapdata_In_Telenav_OSRM_Backend_Docker/)    


## Example By Manual
- [Build Berlin Server with OSM data](./example-berlin-osm.md)

