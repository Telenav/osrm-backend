package trafficcacheindexedbywayid

import (
	"sync"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

type incidentsCache struct {
	m                         sync.RWMutex
	incidents                 map[string]*proxy.Incident
	wayIDBlockedByIncidentIDs map[int64]map[string]struct{} // wayID -> IncidentID,IncidentID,...
}

func newIncidentCache() incidentsCache {
	return incidentsCache{
		sync.RWMutex{},
		map[string]*proxy.Incident{},
		map[int64]map[string]struct{}{},
	}
}

func (i *incidentsCache) clear() {
	i.m.Lock()
	defer i.m.Unlock()

	i.incidents = map[string]*proxy.Incident{}
	i.wayIDBlockedByIncidentIDs = map[int64]map[string]struct{}{}
}

func (i *incidentsCache) isBlockedByIncident(wayID int64) bool {
	i.m.RLock()
	defer i.m.RUnlock()

	if _, ok := i.wayIDBlockedByIncidentIDs[wayID]; ok {
		return true
	}

	return false
}

func (i *incidentsCache) incidentAndAffectedWaysCount() (int, int) {
	i.m.RLock()
	defer i.m.RUnlock()

	return len(i.incidents), len(i.wayIDBlockedByIncidentIDs)
}

func (i *incidentsCache) updateIncidents(incidentResponses []*proxy.IncidentResponse) {
	if len(incidentResponses) == 0 {
		glog.Warning("empty incident responses")
		return
	}

	i.m.Lock()
	defer i.m.Unlock()

	for _, incidentResp := range incidentResponses {
		if incidentResp.Action == proxy.Action_UPDATE || incidentResp.Action == proxy.Action_ADD { //TODO: Action_ADD will be removed soon
			i.unsafeUpdate(incidentResp.Incident)
			continue
		} else if incidentResp.Action == proxy.Action_DELETE {
			i.unsafeDelete(incidentResp.Incident)
			continue
		}

		//undefined
		glog.Errorf("undefined incident action %d, incident %v", incidentResp.Action, incidentResp.Incident)
	}
}

func (i *incidentsCache) unsafeUpdate(incident *proxy.Incident) {
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

	incidentInCache, foundIncidentInCache := i.incidents[incident.IncidentId]
	if foundIncidentInCache {
		i.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
	}
	i.incidents[incident.IncidentId] = incident
	i.unsafeAddWayIDsBlockedByIncidentID(incident.AffectedWayIds, incident.IncidentId)
}

func (i *incidentsCache) unsafeDelete(incident *proxy.Incident) {
	if incident == nil {
		glog.Fatal("empty incident")
		return
	}

	incidentInCache, foundIncidentInCache := i.incidents[incident.IncidentId]
	if foundIncidentInCache {
		i.unsafeDeleteWayIDsBlockedByIncidentID(incidentInCache.AffectedWayIds, incidentInCache.IncidentId)
		delete(i.incidents, incident.IncidentId)
	}
}

func (i *incidentsCache) unsafeDeleteWayIDsBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		if incidentIDs, ok := i.wayIDBlockedByIncidentIDs[wayID]; ok {
			delete(incidentIDs, incidentID)
		}
	}
}

func (i *incidentsCache) unsafeAddWayIDsBlockedByIncidentID(wayIDs []int64, incidentID string) {
	for _, wayID := range wayIDs {
		if incidentIDs, ok := i.wayIDBlockedByIncidentIDs[wayID]; ok {
			incidentIDs[incidentID] = struct{}{} //will do nothing if it's already exist
			continue
		}
		i.wayIDBlockedByIncidentIDs[wayID] = map[string]struct{}{
			incidentID: struct{}{},
		}
	}
}
