package querytraffic

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Inquirer defines interfaces for querying traffic flows and incidents.
type Inquirer interface {
	QueryFlow(int64) *proxy.Flow

	BlockedByIncident(int64) bool
}
