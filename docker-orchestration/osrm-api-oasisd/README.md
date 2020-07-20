# osrm-api-oasisd docker

Image within build `osrm-api-oasisd` binaries, excluding generating pre-processing data by `place-connectivity-gen`, or start web service of `oasisd`.  

## Build Image

- [Dockerfile](./Dockerfile)

```bash
$ cd docker-orchestration/osrm-api-oasisd
$ docker build -t osrm-api-oasisd .   
```
