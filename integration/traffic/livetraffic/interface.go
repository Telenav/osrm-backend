package livetraffic

import (
	"github.com/Telenav/osrm-backend/integration/graph"
	"github.com/Telenav/osrm-backend/integration/traffic"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// Querier defines interfaces for querying traffic flows and incidents.
type Querier traffic.LiveTrafficQuerier

// QuerierByEdge defines interfaces for querying traffic flows and incidents.
type QuerierByEdge interface {
	QueryFlow(graph.Edge) *trafficproxy.Flow
	QueryFlows([]graph.Edge) []*trafficproxy.Flow

	EdgeBlockedByIncident(graph.Edge) bool

	// the second int indicates which Edge blocked by incident.
	EdgesBlockedByIncidents([]graph.Edge) (bool, int)
}
