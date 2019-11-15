package trafficcacheindexedbywayid

import (
	"sync"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Cache is used to cache live traffic and provide query interfaces.
type Cache struct {
	flowsCache     sync.Map
	incidentsCache incidentsCache
}

// New creates a new Cache instance.
func New() *Cache {
	c := Cache{
		sync.Map{},
		newIncidentCache(),
	}
	return &c
}

// Clear all cached traffic flows and incidents.
func (c *Cache) Clear() {
	c.clearFlows()
	c.incidentsCache.clear()
}

// Eat implements trafficproxyclient.Eater inteface.
func (c *Cache) Eat(r proxy.TrafficResponse) {
	c.updateFlows(r.FlowResponses)
	c.incidentsCache.updateIncidents(r.IncidentResponses)
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

	return c.incidentsCache.isBlockedByIncident(wayID)
}

// IncidentCount returns how many incidents in the cache.
// The second returned value is how many ways affected by these incidents.
func (c *Cache) IncidentCount() (int, int) {
	return c.incidentsCache.incidentAndAffectedWaysCount()
}
