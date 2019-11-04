package trafficproxyclient

import (
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
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
	r.MaxSize = 1000
	r.MaxTime = 5
	return &r
}

func (p params) rpcGetTimeout() time.Duration {
	return flags.rpcGetTimeout
}
