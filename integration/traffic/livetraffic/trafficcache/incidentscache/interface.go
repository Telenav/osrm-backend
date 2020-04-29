//Package incidentscache implements cache in memory for blocking-only incidents.
package incidentscache

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/golang/glog"
)

// Cache stores incidents in memory.
type Cache struct {
	m                         sync.RWMutex
	incidents                 map[string]*trafficproxy.Incident
	wayIDBlockedByIncidentIDs map[int64]map[string]struct{} // wayID -> IncidentID,IncidentID,...
}

// New creates a new Cache object to store incidents in memory.
func New() *Cache {
	return &Cache{
		sync.RWMutex{},
		map[string]*trafficproxy.Incident{},
		map[int64]map[string]struct{}{},
	}
}

//Clear clear the cache.
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.incidents = map[string]*trafficproxy.Incident{}
	c.wayIDBlockedByIncidentIDs = map[int64]map[string]struct{}{}
}

// WayBlockedByIncident check whether this wayID is on blocking incident.
func (c *Cache) WayBlockedByIncident(wayID int64) bool {
	c.m.RLock()
	defer c.m.RUnlock()

	if _, ok := c.wayIDBlockedByIncidentIDs[wayID]; ok {
		return true
	}

	return false
}

// Count returns how many incidents in cache.
func (c *Cache) Count() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.incidents)
}

// AffectedWaysCount returns how many ways affected by these incidents in cache.
func (c *Cache) AffectedWaysCount() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.wayIDBlockedByIncidentIDs)
}

// Update updates incidents in cache.
func (c *Cache) Update(incidentResponses []*trafficproxy.IncidentResponse) {
	if len(incidentResponses) == 0 {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, incidentResp := range incidentResponses {
		if incidentResp.Action == trafficproxy.Action_UPDATE {
			c.unsafeUpdate(incidentResp.Incident)
			continue
		} else if incidentResp.Action == trafficproxy.Action_DELETE {
			c.unsafeDelete(incidentResp.Incident)
			continue
		}

		//undefined
		glog.Errorf("undefined incident action %d, incident %v", incidentResp.Action, incidentResp.Incident)
	}
}
