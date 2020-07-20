FROM telenavmap/osrm-backend:latest as builder

COPY ./docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

EXPOSE 8090
ENTRYPOINT ["/docker-entrypoint.sh"]
