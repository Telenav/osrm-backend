

FROM telenavmap/osrm-backend-dev:v1 as builder

RUN mkdir -p /osrm-build
WORKDIR /workspace/

# specify a branch/tag/commit to checkout, master by default
ARG GIT_COMMIT=master  
RUN echo "Building on ${GIT_COMMIT}" && \
  git clone -n https://github.com/Telenav/osrm-backend.git && \
  cd osrm-backend && \
  git checkout ${GIT_COMMIT}
  
# Alternatively `-- -j` to enable parallel compiling.
# It can not be enabled by default due to GitHub Actions runner's resource limitation.
ARG CMAKE_BUILD_EXTRA_ARGS=
RUN mkdir -p osrm-backend/build && \
  cd osrm-backend/build && \
  cmake .. -DENABLE_LTO=On && \
  cmake --build . ${CMAKE_BUILD_EXTRA_ARGS} && \
  cp /workspace/osrm-backend/build/osrm* /osrm-build/ && \
  cp -r /workspace/osrm-backend/profiles /osrm-build/

RUN cd /workspace/osrm-backend/integration && \
  go test -short -v ./... && \
  go install ./... && \
  ls -lh /workspace/go/bin


FROM telenavmap/osrm-backend-dev:v1 as runstage
#FROM debian:stretch-slim as runstage
#RUN DEBIAN_FRONTEND=noninteractive apt-get update && \
#  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends libboost-program-options1.62.0 libboost-regex1.62.0 \
#    libboost-date-time1.62.0 libboost-chrono1.62.0 libboost-filesystem1.62.0 \
#    libboost-iostreams1.62.0 libboost-thread1.62.0 expat liblua5.2-0 libtbb2 curl ca-certificates && \
#  rm -rf /var/lib/apt/lists/*

RUN mkdir -p /osrm-build /osrm-data /osrm-logs

COPY --from=builder /osrm-build /osrm-build/
COPY --from=builder /workspace/go/bin /osrm-build/
COPY ./docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

# for shown on runtime
ARG IMAGE_TAG
ENV IMAGE_TAG ${IMAGE_TAG}

EXPOSE 5000
ENTRYPOINT ["/docker-entrypoint.sh"]

