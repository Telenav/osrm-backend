package trafficproxyclient

import (
	"testing"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
)

func TestIsBlockingFlow(t *testing.T) {
	cases := []struct {
		flow   *proxy.Flow
		expect bool
	}{
		{flow: nil, expect: false},
		{flow: &proxy.Flow{Speed: 2.0}, expect: false},
		{flow: &proxy.Flow{Speed: 2.0, TrafficLevel: proxy.TrafficLevel_SLOW_SPEED}, expect: false},
		{flow: &proxy.Flow{Speed: 1.0}, expect: false},
		{flow: &proxy.Flow{TrafficLevel: proxy.TrafficLevel_CLOSED}, expect: true},
		{flow: &proxy.Flow{Speed: 0.9}, expect: true},
		{flow: &proxy.Flow{TrafficLevel: proxy.TrafficLevel_CLOSED, Speed: 1.0}, expect: true},
		{flow: &proxy.Flow{TrafficLevel: proxy.TrafficLevel_FREE_FLOW, Speed: 0.9}, expect: true},
		{flow: &proxy.Flow{TrafficLevel: proxy.TrafficLevel_CLOSED, Speed: 0}, expect: true},
	}

	for _, c := range cases {
		result := IsBlockingFlow(c.flow)
		if result != c.expect {
			t.Errorf("flow: %v, expect: %t, but IsBlockingFlow return %t", c.flow, c.expect, result)
		}
	}

}

func TestIsBlockingIncident(t *testing.T) {

	cases := []struct {
		incident *proxy.Incident
		expect   bool
	}{
		{incident: nil, expect: false},
		{incident: &proxy.Incident{IsBlocking: false}, expect: false},
		{incident: &proxy.Incident{IsBlocking: true}, expect: true},
	}

	for _, c := range cases {
		result := IsBlockingIncident(c.incident)
		if result != c.expect {
			t.Errorf("incident: %v, expect: %t, but IsBlockingIncident return %t", c.incident, c.expect, result)
		}
	}
}
