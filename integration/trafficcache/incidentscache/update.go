package incidentscache

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

func (c *Cache) unsafeUpdate(incident *proxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}
	if len(incident.AffectedWayIds) == 0 {
		glog.Warningf("empty AffectedWayIds in incident %v", incident)
		return
	}
	if !incident.IsBlocking {
		return // we only take care of blocking incidents
	}

	incidentInCache, foundIncidentInCache := c.incidents[incident.IncidentId]
	if foundIncidentInCache {
		c.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
			c.unsafeDeleteEdgesBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		}
	}
	c.incidents[incident.IncidentId] = incident
	c.unsafeAddWayIDsBlockedByIncidentID(incident.AffectedWayIds, incident.IncidentId)
	if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
		c.unsafeAddEdgesBlockedByIncidentID(incident.AffectedWayIds, incident.IncidentId)
	}
}

func (c *Cache) unsafeDelete(incident *proxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}

	incidentInCache, foundIncidentInCache := c.incidents[incident.IncidentId]
	if foundIncidentInCache {
		c.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
			c.unsafeDeleteEdgesBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		}
		delete(c.incidents, incident.IncidentId)
	}
}

func (c *Cache) unsafeDeleteWayIDsBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		if incidentIDs, ok := c.wayIDBlockedByIncidentIDs[wayID]; ok {
			delete(incidentIDs, incidentID)
			if len(incidentIDs) == 0 { // the wayID doesn't blocked by incident anymore
				delete(c.wayIDBlockedByIncidentIDs, wayID)
			}
		}
	}
}

func (c *Cache) unsafeAddWayIDsBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		if incidentIDs, ok := c.wayIDBlockedByIncidentIDs[wayID]; ok {
			incidentIDs[incidentID] = struct{}{} //will do nothing if it's already exist
			continue
		}
		c.wayIDBlockedByIncidentIDs[wayID] = map[string]struct{}{
			incidentID: struct{}{},
		}
	}
}

func (c *Cache) unsafeDeleteEdgesBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		edges := c.wayID2Edges.WayID2Edges(wayID)

		for _, edge := range edges {
			if incidentIDs, ok := c.edgeBlockedByIncidentIDs[edge]; ok {
				delete(incidentIDs, incidentID)
				if len(incidentIDs) == 0 { // the edge doesn't blocked by incident anymore
					delete(c.edgeBlockedByIncidentIDs, edge)
				}
			}
		}
	}
}

func (c *Cache) unsafeAddEdgesBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		edges := c.wayID2Edges.WayID2Edges(wayID)

		for _, edge := range edges {
			if incidentIDs, ok := c.edgeBlockedByIncidentIDs[edge]; ok {
				incidentIDs[incidentID] = struct{}{} //will do nothing if it's already exist
				continue
			}
			c.edgeBlockedByIncidentIDs[edge] = map[string]struct{}{
				incidentID: struct{}{},
			}
		}
	}
}
