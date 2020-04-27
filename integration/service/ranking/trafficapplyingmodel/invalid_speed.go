package trafficapplyingmodel

import (
	"github.com/Telenav/osrm-backend/integration/util"
)

// Invalid live/historical traffic speeds for applying.
const (
	InvalidTrafficSpeed        = -1
	InvalidTrafficSpeedFloat64 = float64(InvalidTrafficSpeed)
)

// IsInvalidSpeed decide whether the speed is valid(>=0) or not.
func IsInvalidSpeed(speed float64) bool {
	return util.FloatEquals(speed, InvalidTrafficSpeedFloat64) || speed < 0
}
