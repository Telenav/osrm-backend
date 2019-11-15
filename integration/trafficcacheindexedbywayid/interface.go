package trafficcacheindexedbywayid

import (
	"sync"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/trafficcache/incidentscache"
	"github.com/golang/glog"
)

// Cache is used to cache live traffic and provide query interfaces.
type Cache struct {
	flowsCache     sync.Map
	incidentsCache incidentscache.Cache
}

// New creates a new Cache instance.
func New() *Cache {
	c := Cache{
		sync.Map{},
		incidentscache.New(),
	}
	return &c
}

// Clear all cached traffic flows and incidents.
func (c *Cache) Clear() {
	c.clearFlows()
	c.incidentsCache.Clear()
}

// Eat implements trafficproxyclient.Eater inteface.
func (c *Cache) Eat(r proxy.TrafficResponse) {
	glog.V(1).Infof("new traffic for cache, flows: %d, incidents: %d", len(r.FlowResponses), len(r.IncidentResponses))
	c.updateFlows(r.FlowResponses)
	c.incidentsCache.Update(r.IncidentResponses)
}

// QueryFlow returns Live Traffic Flow if exist.
func (c *Cache) QueryFlow(wayID int64) *proxy.Flow {

	return c.queryFlow(wayID)
}

// FlowCount returns how many flows in the cache.
func (c *Cache) FlowCount() int64 {
	return c.flowCount()
}

// IsBlockedByIncident checks whether the way has blocking incident.
func (c *Cache) IsBlockedByIncident(wayID int64) bool {

	return c.incidentsCache.IsWayBlockedByIncident(wayID)
}

// IncidentCount returns how many incidents in the cache.
func (c *Cache) IncidentCount() (int, int) {
	return c.incidentsCache.IncidentsAndAffectedWaysCount()
}

// AffectedWaysCount returns how many ways affected by these incidents in cache.
func (c *Cache) AffectedWaysCount() int {
	return c.incidentsCache.AffectedWaysCount()
}

// IncidentsAndAffectedWaysCount returns how many incidents in cache and how many ways affected by these incidents.
func (c *Cache) IncidentsAndAffectedWaysCount() (int, int) {
	return c.incidentsCache.IncidentsAndAffectedWaysCount()
}
