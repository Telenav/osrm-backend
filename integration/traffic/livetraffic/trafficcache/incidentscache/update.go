package incidentscache

import (
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/golang/glog"
)

func (c *Cache) unsafeUpdate(incident *trafficproxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}
	if len(incident.AffectedWayIDs) == 0 {
		glog.Warningf("empty AffectedWayIds in incident %v", incident)
		return
	}
	if !incident.IsBlocking {
		return // we only take care of blocking incidents
	}

	incidentInCache, foundIncidentInCache := c.incidents[incident.IncidentID]
	if foundIncidentInCache {
		if incidentInCache.Timestamp > incident.Timestamp {
			return // ignore older incident
		}
		c.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIDs, incidentInCache.IncidentID)
		if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
			c.unsafeDeleteEdgesBlockedByIncidentID(incidentInCache.AffectedWayIDs, incidentInCache.IncidentID)
		}
	}
	c.incidents[incident.IncidentID] = incident
	c.unsafeAddWayIDsBlockedByIncidentID(incident.AffectedWayIDs, incident.IncidentID)
	if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
		c.unsafeAddEdgesBlockedByIncidentID(incident.AffectedWayIDs, incident.IncidentID)
	}
}

func (c *Cache) unsafeDelete(incident *trafficproxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}

	incidentInCache, foundIncidentInCache := c.incidents[incident.IncidentID]
	if foundIncidentInCache {
		c.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIDs, incidentInCache.IncidentID)
		if c.wayID2Edges != nil && c.edgeBlockedByIncidentIDs != nil {
			c.unsafeDeleteEdgesBlockedByIncidentID(incidentInCache.AffectedWayIDs, incidentInCache.IncidentID)
		}
		delete(c.incidents, incident.IncidentID)
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
