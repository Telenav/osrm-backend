package trafficcacheindexedbyedge

import (
	"github.com/Telenav/osrm-backend/integration/graph"
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/trafficcache/flowscacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/trafficcache/incidentscache"
	"github.com/golang/glog"
)

// Cache is used to cache live traffic and provide query interfaces.
type Cache struct {
	Flows     *flowscacheindexedbyedge.Cache
	Incidents *incidentscache.Cache
}

// New creates a new Cache instance.
func New(wayID2Edges graph.WayID2EdgesMapping) *Cache {
	c := Cache{
		flowscacheindexedbyedge.New(wayID2Edges),
		incidentscache.NewWithEdgeIndexing(wayID2Edges),
	}
	return &c
}

// Clear all cached traffic flows and incidents.
func (c *Cache) Clear() {
	c.Flows.Clear()
	c.Incidents.Clear()
}

// Eat implements trafficeater.Eater inteface.
func (c *Cache) Eat(r proxy.TrafficResponse) {
	glog.V(1).Infof("new traffic for cache, flows: %d, incidents: %d", len(r.FlowResponses), len(r.IncidentResponses))
	c.Flows.Update(r.FlowResponses)
	c.Incidents.Update(r.IncidentResponses)
}

// QueryFlow returns Live Traffic Flow if exist.
func (c *Cache) QueryFlow(e graph.Edge) *proxy.Flow {
	return c.Flows.QueryEdge(e)
}

// QueryFlows returns Live Traffic Flows if exist.
func (c *Cache) QueryFlows(e []graph.Edge) []*proxy.Flow {
	return c.Flows.QueryEdges(e)
}

// EdgeBlockedByIncident check whether this Edge is on blocking incident.
func (c *Cache) EdgeBlockedByIncident(e graph.Edge) bool {
	return c.Incidents.EdgeBlockedByIncident(e)
}

// EdgesBlockedByIncidents check whether this Edge is on blocking incidents.
// the second return indicates the blocked edge index of input array if exist.
func (c *Cache) EdgesBlockedByIncidents(e []graph.Edge) (bool, int) {
	return c.Incidents.EdgesBlockedByIncidents(e)
}
