package traffic

import "time"

// HistoricalSpeedQuerier provide interfaces for Querying Historical Speeds
type HistoricalSpeedQuerier interface {

	// QueryHistoricalSpeed return the speed for a way at a specified time
	// wayID: positive means travel forward, which is travel following edge's point sequence, negative means backward
	// t: UTC time
	QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool)
}
