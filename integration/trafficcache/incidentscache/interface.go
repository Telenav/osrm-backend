//Package incidentscache implements cache in memory for blocking-only incidents.
package incidentscache

import (
	"sync"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// Cache stores incidents in memory.
type Cache struct {
	m                         sync.RWMutex
	incidents                 map[string]*proxy.Incident
	wayIDBlockedByIncidentIDs map[int64]map[string]struct{} // wayID -> IncidentID,IncidentID,...
}

// New creates a new Cache object to store incidents in memory.
func New() Cache {
	return Cache{
		sync.RWMutex{},
		map[string]*proxy.Incident{},
		map[int64]map[string]struct{}{},
	}
}

//Clear clear the cache.
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.incidents = map[string]*proxy.Incident{}
	c.wayIDBlockedByIncidentIDs = map[int64]map[string]struct{}{}
}

// IsWayBlockedByIncident check whether this wayID is on blocking incident.
func (c *Cache) IsWayBlockedByIncident(wayID int64) bool {
	c.m.RLock()
	defer c.m.RUnlock()

	if _, ok := c.wayIDBlockedByIncidentIDs[wayID]; ok {
		return true
	}

	return false
}

// IncidentCount returns how many incidents in cache.
func (c *Cache) IncidentCount() int {
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

// IncidentsAndAffectedWaysCount returns how many incidents in cache and how many ways affected by these incidents.
func (c *Cache) IncidentsAndAffectedWaysCount() (int, int) {
	c.m.RLock()
	defer c.m.RUnlock()

	return len(c.incidents), len(c.wayIDBlockedByIncidentIDs)
}

// Update updates incidents in cache.
func (c *Cache) Update(incidentResponses []*proxy.IncidentResponse) {
	if len(incidentResponses) == 0 {
		return
	}

	c.m.Lock()
	defer c.m.Unlock()

	for _, incidentResp := range incidentResponses {
		if incidentResp.Action == proxy.Action_UPDATE || incidentResp.Action == proxy.Action_ADD { //TODO: Action_ADD will be removed soon
			c.unsafeUpdate(incidentResp.Incident)
			continue
		} else if incidentResp.Action == proxy.Action_DELETE {
			c.unsafeDelete(incidentResp.Incident)
			continue
		}

		//undefined
		glog.Errorf("undefined incident action %d, incident %v", incidentResp.Action, incidentResp.Incident)
	}
}
