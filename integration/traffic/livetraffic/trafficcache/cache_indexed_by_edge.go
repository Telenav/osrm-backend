package trafficcache

import (
	"github.com/Telenav/osrm-backend/integration/graph"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/flowscacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/incidentscache"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/wayidsmap"
	"github.com/golang/glog"
)

// CacheIndexedByEdge is used to cache live traffic and provide query interfaces.
type CacheIndexedByEdge struct {
	Flows     *flowscacheindexedbyedge.Cache
	Incidents *incidentscache.Cache
}

// NewCacheIndexedByEdge creates a new CacheIndexedByEdge instance.
func NewCacheIndexedByEdge(wayID2Edges wayidsmap.Way2Edges) *CacheIndexedByEdge {
	c := CacheIndexedByEdge{
		flowscacheindexedbyedge.New(wayID2Edges),
		incidentscache.NewWithEdgeIndexing(wayID2Edges),
	}
	return &c
}

// Clear all cached traffic flows and incidents.
func (c *CacheIndexedByEdge) Clear() {
	c.Flows.Clear()
	c.Incidents.Clear()
}

// Eat implements livetraffic.Eater inteface.
func (c *CacheIndexedByEdge) Eat(r trafficproxy.TrafficResponse) {
	glog.V(1).Infof("new traffic for cache, flows: %d, incidents: %d", len(r.FlowResponses), len(r.IncidentResponses))
	c.Flows.Update(r.FlowResponses)
	c.Incidents.Update(r.IncidentResponses)
}

// QueryFlow returns Live Traffic Flow if exist.
func (c *CacheIndexedByEdge) QueryFlow(e graph.Edge) *trafficproxy.Flow {
	return c.Flows.QueryByEdge(e)
}

// QueryFlows returns Live Traffic Flows if exist.
func (c *CacheIndexedByEdge) QueryFlows(e []graph.Edge) []*trafficproxy.Flow {
	return c.Flows.QueryByEdges(e)
}

// EdgeBlockedByIncident check whether this Edge is on blocking incident.
func (c *CacheIndexedByEdge) EdgeBlockedByIncident(e graph.Edge) bool {
	return c.Incidents.EdgeBlockedByIncident(e)
}

// EdgesBlockedByIncidents check whether this Edge is on blocking incidents.
// the second return indicates the blocked edge index of input array if exist.
func (c *CacheIndexedByEdge) EdgesBlockedByIncidents(e []graph.Edge) (bool, int) {
	return c.Incidents.EdgesBlockedByIncidents(e)
}
