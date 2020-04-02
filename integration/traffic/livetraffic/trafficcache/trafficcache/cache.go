package trafficcache

import (
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/flowscache"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/incidentscache"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/golang/glog"
)

// Cache is used to cache live traffic and provide query interfaces.
type Cache struct {
	Flows     *flowscache.Cache
	Incidents *incidentscache.Cache
}

// New creates a new Cache instance.
func New() *Cache {
	c := Cache{
		flowscache.New(),
		incidentscache.New(),
	}
	return &c
}

// Clear all cached traffic flows and incidents.
func (c *Cache) Clear() {
	c.Flows.Clear()
	c.Incidents.Clear()
}

// Eat implements trafficeater.Eater inteface.
func (c *Cache) Eat(r trafficproxy.TrafficResponse) {
	glog.V(1).Infof("new traffic for cache, flows: %d, incidents: %d", len(r.FlowResponses), len(r.IncidentResponses))
	c.Flows.Update(r.FlowResponses)
	c.Incidents.Update(r.IncidentResponses)
}

// QueryFlow returns Live Traffic Flow if exist.
func (c *Cache) QueryFlow(wayID int64) *trafficproxy.Flow {
	return c.Flows.Query(wayID)
}

// BlockedByIncident check whether this wayID is on blocking incident.
func (c *Cache) BlockedByIncident(wayID int64) bool {
	return c.Incidents.WayBlockedByIncident(wayID)
}
