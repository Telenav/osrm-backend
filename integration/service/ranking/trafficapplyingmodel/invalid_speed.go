package trafficapplyingmodel

import (
	"github.com/Telenav/osrm-backend/integration/util"
)

// Invalid live/historical traffic speeds for applying.
const (
	InvalidSpeed            = -1
	InvalidLiveTrafficSpeed = float32(InvalidSpeed) // Speed is float32 in trafficproxy.Flow
	InvalidHistoricalSpeed  = float64(InvalidSpeed)
)

// IsInvalidSpeed decide whether the speed is valid(>=0) or not.
func IsInvalidSpeed(speed float64) bool {
	return util.Float32Equal(float32(speed), float32(InvalidSpeed)) || speed < 0
}
