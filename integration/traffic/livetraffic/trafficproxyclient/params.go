package trafficproxyclient

import (
	"strings"
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// params is used to group request parameters together.
type params struct{}

func (p params) newTrafficSource() *trafficproxy.TrafficSource {
	t := trafficproxy.TrafficSource{}
	t.Region = flags.region
	t.TrafficProvider = flags.trafficProvider
	t.MapProvider = flags.mapProvider

	t.Subregion = []string{}
	if len(flags.subregion) > 0 {
		subregionStr := strings.TrimSuffix(flags.subregion, ",") // trim last ',' if exist
		t.Subregion = strings.Split(subregionStr, ",")
	}
	return &t
}

func (p params) newTrafficType() []trafficproxy.TrafficType {
	t := []trafficproxy.TrafficType{}
	if flags.flow {
		t = append(t, trafficproxy.TrafficType_FLOW)
	}
	if flags.incident {
		t = append(t, trafficproxy.TrafficType_INCIDENT)
	}
	return t
}

func (p params) newStreamingRule() *trafficproxy.TrafficStreamingDeltaRequest_StreamingRule {
	var r trafficproxy.TrafficStreamingDeltaRequest_StreamingRule
	r.MaxSize = int32(flags.streamingDeltaMaxSize)
	r.MaxTime = int32(flags.streamingDeltaMaxTime.Seconds())
	return &r
}

func (p params) rpcGetTimeout() time.Duration {
	return flags.rpcGetTimeout
}
