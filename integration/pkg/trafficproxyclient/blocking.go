package trafficproxyclient

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
)

const blockingSpeedThreshold = 1 // Think it's blocking if flow speed smaller than this threshold.

// IsBlockingFlow tests whether the Flow is blocking or not.
func IsBlockingFlow(f *proxy.Flow) bool {
	if f == nil {
		return false
	}

	return f.TrafficLevel == proxy.TrafficLevel_CLOSED || f.Speed < blockingSpeedThreshold
}

// IsBlockingIncident tests whether the incident is blocking or not.
func IsBlockingIncident(incident *proxy.Incident) bool {
	if incident == nil {
		return false
	}

	return incident.IsBlocking
}
