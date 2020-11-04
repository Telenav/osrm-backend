package mock

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/util/speedunit"
)

// FixedTraffic implements both traffic.HistoricalSpeedQuerier and traffic.LiveTrafficQuerier but always return fixed speed/level.
type FixedTraffic struct {
	fixedSpeed        float64 // km/h
	fixedTrafficLevel trafficproxy.TrafficLevel
}

// NewFixedTraffic creates a new fixed traffic mock object.The input speed unit is km/h.
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
		Speed:        float32(speedunit.ConvertKPH2MPS(f.fixedSpeed)), // km/h -> m/s since traffic flow's speed unit is m/s.
		TrafficLevel: f.fixedTrafficLevel,
		Timestamp:    1588840770,
	}
}

// QueryHistoricalSpeed implements QueryHistoricalSpeed of traffic.HistoricalSpeedQuerier but always return fixed speed.
func (f FixedTraffic) QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool) {
	return f.fixedSpeed, true
}
