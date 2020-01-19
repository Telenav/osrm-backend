//Package incidentscache implements cache in memory for blocking-only incidents.
package incidentscache

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/graph"
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/wayidsmap"
	"github.com/golang/glog"
)

// Cache stores incidents in memory.
type Cache struct {
	m                         sync.RWMutex
	incidents                 map[string]*trafficproxy.Incident
	wayIDBlockedByIncidentIDs map[int64]map[string]struct{} // wayID -> IncidentID,IncidentID,...

	// optional
	wayID2Edges              wayidsmap.Way2Edges
	edgeBlockedByIncidentIDs map[graph.Edge]map[string]struct{} // edge -> IncidentID,IncidentID,...
}

// New creates a new Cache object to store incidents in memory.
func New() *Cache {
	return &Cache{
		sync.RWMutex{},
		map[string]*trafficproxy.Incident{},
		map[int64]map[string]struct{}{},
		nil,
		nil,
	}
}

// NewWithEdgeIndexing creates a new Cache object to store incidents in memory, with also Edge indexing support.
func NewWithEdgeIndexing(wayID2Edges wayidsmap.Way2Edges) *Cache {
	if wayID2Edges == nil {
		glog.Fatal("empty wayID2Edges")
		return nil
	}

	return &Cache{
		sync.RWMutex{},
		map[string]*trafficproxy.Incident{},
		map[int64]map[string]struct{}{},
		wayID2Edges,
		map[graph.Edge]map[string]struct{}{},
	}
}

//Clear clear the cache.
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.incidents = map[string]*trafficproxy.Incident{}
	c.wayIDBlockedByIncidentIDs = map[int64]map[string]struct{}{}
	if c.edgeBlockedByIncidentIDs != nil {
		c.edgeBlockedByIncidentIDs = map[graph.Edge]map[string]struct{}{}
	}
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

// EdgeBlockedByIncident check whether this Edge is on blocking incident.
func (c *Cache) EdgeBlockedByIncident(edge graph.Edge) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.edgeBlockedByIncidentIDs == nil {
		return false
	}

	if _, ok := c.edgeBlockedByIncidentIDs[edge]; ok {
		return true
	}

	return false
}

// EdgesBlockedByIncidents check whether this Edge is on blocking incidents.
// the second return indicates the blocked edge index of input array if exist.
func (c *Cache) EdgesBlockedByIncidents(edges []graph.Edge) (bool, int) {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.edgeBlockedByIncidentIDs == nil {
		return false, -1
	}

	for i := range edges {
		if _, ok := c.edgeBlockedByIncidentIDs[edges[i]]; ok {
			return true, i
		}
	}

	return false, -1
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

// AffectedEdgesCount returns how many edges affected by these incidents in cache.
func (c *Cache) AffectedEdgesCount() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.edgeBlockedByIncidentIDs)
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
