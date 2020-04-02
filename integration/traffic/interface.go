package traffic

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// HistoricalSpeedQuerier provide interfaces for Querying Historical Speeds
type HistoricalSpeedQuerier interface {

	// QueryHistoricalSpeed return the speed for a way at a specified time
	// wayID: positive means travel forward, which is travel following edge's point sequence, negative means backward
	// t: UTC time
	QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool)
}

// LiveTrafficQuerier defines interfaces for querying traffic flows and incidents.
type LiveTrafficQuerier interface {

	// QueryFlow return the live traffic flow for a way
	// wayID: positive means travel forward, which is travel following edge's point sequence, negative means backward
	QueryFlow(wayID int64) *trafficproxy.Flow

	// BlockedByIncident returns whether a way blocked by live traffic incident
	// wayID: positive means travel forward, which is travel following edge's point sequence, negative means backward
	BlockedByIncident(wayID int64) bool
}
