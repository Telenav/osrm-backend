

ARG FROM_TAG=latest
FROM telenavmap/osrm-backend:${FROM_TAG} 


COPY ./map/map.osrm.cell_metrics /osrm-data/
COPY ./map/map.osrm.cells /osrm-data/
COPY ./map/map.osrm.cnbg /osrm-data/
COPY ./map/map.osrm.cnbg_to_ebg /osrm-data/
COPY ./map/map.osrm.datasource_names /osrm-data/
COPY ./map/map.osrm.ebg /osrm-data/
COPY ./map/map.osrm.ebg_nodes /osrm-data/
COPY ./map/map.osrm.edges /osrm-data/
COPY ./map/map.osrm.enw /osrm-data/
COPY ./map/map.osrm.fileIndex /osrm-data/
COPY ./map/map.osrm.geometry /osrm-data/
COPY ./map/map.osrm.icd /osrm-data/
COPY ./map/map.osrm.maneuver_overrides /osrm-data/
COPY ./map/map.osrm.mldgr /osrm-data/
COPY ./map/map.osrm.names /osrm-data/
COPY ./map/map.osrm.nbg_nodes /osrm-data/
COPY ./map/map.osrm.partition /osrm-data/
COPY ./map/map.osrm.properties /osrm-data/
COPY ./map/map.osrm.ramIndex /osrm-data/
COPY ./map/map.osrm.restrictions /osrm-data/
COPY ./map/map.osrm.timestamp /osrm-data/
COPY ./map/map.osrm.tld /osrm-data/
COPY ./map/map.osrm.tls /osrm-data/
COPY ./map/map.osrm.turn_duration_penalties /osrm-data/
COPY ./map/map.osrm.turn_penalties_index /osrm-data/
COPY ./map/map.osrm.turn_weight_penalties /osrm-data/
COPY ./map/wayid2nodeids.csv.snappy /osrm-data/

EXPOSE 5000
ENTRYPOINT ["/docker-entrypoint.sh"]

