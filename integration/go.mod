module github.com/Telenav/osrm-backend/integration

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/qedus/osmpbf v1.1.0
	google.golang.org/grpc v1.22.0
)

replace github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy => ./pkg/gen-trafficproxy
