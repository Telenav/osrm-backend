package querytraffic

import (
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Inquirer defines interfaces for querying traffic flows and incidents.
type Inquirer interface {
	QueryFlow(int64) *trafficproxy.Flow

	BlockedByIncident(int64) bool
}
