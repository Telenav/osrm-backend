package trafficproxyclient

import proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"

func newTrafficSource() *proxy.TrafficSource {
	t := proxy.TrafficSource{}
	t.Region = flags.region
	t.TrafficProvider = flags.trafficProvider
	t.MapProvider = flags.mapProvider
	return &t
}

func newTrafficType() []proxy.TrafficType {
	t := []proxy.TrafficType{}
	if flags.flow {
		t = append(t, proxy.TrafficType_FLOW)
	}
	if flags.incident {
		t = append(t, proxy.TrafficType_INCIDENT)
	}
	return t
}
