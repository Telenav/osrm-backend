package trafficproxyclient

import proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"

func newTrafficSource() *proxy.TrafficSource {
	t := proxy.TrafficSource{}
	t.Region = flags.Region
	t.TrafficProvider = flags.TrafficProvider
	t.MapProvider = flags.MapProvider
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
