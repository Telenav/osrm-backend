package inquirerbyedge

import (
	"github.com/Telenav/osrm-backend/integration/graph"
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Inquirer defines interfaces for querying traffic flows and incidents.
type Inquirer interface {
	QueryFlow(graph.Edge) *proxy.Flow
	QueryFlows([]graph.Edge) []*proxy.Flow

	EdgeBlockedByIncident(graph.Edge) bool

	// the second int indicates which Edge blocked by incident.
	EdgesBlockedByIncidents([]graph.Edge) (bool, int)
}
