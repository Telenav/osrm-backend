package mock

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// FixedTraffic implements both traffic.HistoricalSpeedQuerier and traffic.LiveTrafficQuerier but always return fixed speed/level.
type FixedTraffic struct {
	fixedSpeed        float64
	fixedTrafficLevel trafficproxy.TrafficLevel
}

// NewFixedTraffic creates a new fixed traffic mock object.
func NewFixedTraffic(fixedSpeed float64, fixedTrafficLevel trafficproxy.TrafficLevel) FixedTraffic {
	return FixedTraffic{fixedSpeed, fixedTrafficLevel}
}

// BlockedByIncident implements BlockedByIncident of traffic.LiveTrafficQuerier but always return false.
func (f FixedTraffic) BlockedByIncident(wayID int64) bool {
	return false
}

// QueryFlow implements QueryFlow of traffic.LiveTrafficQuerier but always return fixed speed.
func (f FixedTraffic) QueryFlow(wayID int64) *trafficproxy.Flow {
	return &trafficproxy.Flow{
		WayID:        wayID,
		Speed:        float32(f.fixedSpeed),
		TrafficLevel: f.fixedTrafficLevel,
		Timestamp:    1588840770,
	}
}

// QueryHistoricalSpeed implements QueryHistoricalSpeed of traffic.HistoricalSpeedQuerier but always return fixed speed.
func (f FixedTraffic) QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool) {
	return f.fixedSpeed, true
}
