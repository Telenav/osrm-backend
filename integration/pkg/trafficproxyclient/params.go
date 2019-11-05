package trafficproxyclient

import (
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// params is used to group request parameters together.
type params struct{}

func (p params) newTrafficSource() *proxy.TrafficSource {
	t := proxy.TrafficSource{}
	t.Region = flags.region
	t.TrafficProvider = flags.trafficProvider
	t.MapProvider = flags.mapProvider
	return &t
}

func (p params) newTrafficType() []proxy.TrafficType {
	t := []proxy.TrafficType{}
	if flags.flow {
		t = append(t, proxy.TrafficType_FLOW)
	}
	if flags.incident {
		t = append(t, proxy.TrafficType_INCIDENT)
	}
	return t
}

func (p params) newStreamingRule() *proxy.TrafficStreamingDeltaRequest_StreamingRule {
	var r proxy.TrafficStreamingDeltaRequest_StreamingRule
	r.MaxSize = int32(flags.streamingDeltaMaxSize)
	r.MaxTime = int32(flags.streamingDeltaMaxTime.Seconds())
	return &r
}

func (p params) rpcGetTimeout() time.Duration {
	return flags.rpcGetTimeout
}
