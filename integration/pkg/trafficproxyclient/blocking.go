package trafficproxyclient

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
	"github.com/golang/glog"
)

const blockingSpeedThreshold = 1 // Think it's blocking if flow speed smaller than this threshold.

// IsBlockingFlow tests whether the Flow is blocking or not.
func IsBlockingFlow(f *proxy.Flow) bool {
	if f == nil {
		glog.Fatal("empty flow")
	}

	return f.TrafficLevel == proxy.TrafficLevel_CLOSED || f.Speed < blockingSpeedThreshold
}

// IsBlockingIncident tests whether the incident is blocking or not.
func IsBlockingIncident(incident *proxy.Incident) bool {
	if incident == nil {
		glog.Fatal("empty incident")
	}

	return incident.IsBlocking
}
