package mock

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// EmptyTraffic implements both traffic.HistoricalSpeedQuerier and traffic.LiveTrafficQuerier but no any data.
type EmptyTraffic struct{}

// BlockedByIncident implements BlockedByIncident of traffic.LiveTrafficQuerier but always return false.
func (e EmptyTraffic) BlockedByIncident(wayID int64) bool {
	return false
}

// QueryFlow implements QueryFlow of traffic.LiveTrafficQuerier but always return nil.
func (e EmptyTraffic) QueryFlow(wayID int64) *trafficproxy.Flow {
	return nil
}

// QueryHistoricalSpeed implements QueryHistoricalSpeed of traffic.HistoricalSpeedQuerier but always return false.
func (e EmptyTraffic) QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool) {
	return 0, false
}
