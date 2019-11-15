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
	}
	c.incidents[incident.IncidentId] = incident
	c.unsafeAddWayIDsBlockedByIncidentID(incident.AffectedWayIds, incident.IncidentId)
}

func (c *Cache) unsafeDelete(incident *proxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}

	incidentInCache, foundIncidentInCache := c.incidents[incident.IncidentId]
	if foundIncidentInCache {
		c.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		delete(c.incidents, incident.IncidentId)
	}
}

func (c *Cache) unsafeDeleteWayIDsBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		if incidentIDs, ok := c.wayIDBlockedByIncidentIDs[wayID]; ok {
			delete(incidentIDs, incidentID)
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
