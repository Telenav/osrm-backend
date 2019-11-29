# Docker Orchestration

## Docker Images 
### osrm-backend-dev
Base image for telenav osrm-backend development, include all building and running dependencies.     
```bash
$ docker pull telenavmap/osrm-backend-dev 
```

### osrm-backend
Image within built osrm binaries(`osrm-extract/osrm-partition/osrm-customize/...`) and running dependencies. It can be used to **compile data** or **startup routed**.      

```bash
$ docker push telenavmap/osrm-backend 
```

See details in [osrm-backend docker](./osrm-backend/)

### osrm-backend-within-mapdata
Image based on [osrm-backend docker image](#osrm-backend) and put compiled mapdata inside.          
NOTE: It's a temporary workaround for easily run in k8s. It's NOT a good idea to put mapdata in image directly since the map data is too big. For long-term, discussing in [#93](https://github.com/Telenav/osrm-backend/issues/93).      

### osrm-frontend
Image contains web tool to check routing and guidance result.  
It uses MapBox GL JS and apply routing response on top of Mapbox vector tiles.  
See details in [osrm-frontend-docker](./osrm-frontend-docker/README.md)

## Kubernetes Deployment
### k8s-rolling-update
Use kubernetes rolling update deployment strategy for timed replace container with new one. Latest traffic will be used during container startup.  
See details in [k8s rolling update](./k8s-rolling-update/)
