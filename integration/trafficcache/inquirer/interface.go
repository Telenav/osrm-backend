package inquirer

import (
	"github.com/Telenav/osrm-backend/integration/graph"
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Inquirer defines interfaces for querying traffic flows and incidents.
type Inquirer interface {
	QueryFlow(graph.WayID) *proxy.Flow

	BlockedByIncident(graph.WayID) bool
}
