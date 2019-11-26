package incidentscache

import (
	"testing"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

func TestIncidentsCache(t *testing.T) {
	presetIncidents := []*proxy.Incident{
		&proxy.Incident{
			IncidentId:            "TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1",
			AffectedWayIds:        []int64{100663296, -1204020275, 100663296, -1204020274, 100663296, -916744017, 100663296, -1204020245, 100663296, -1194204646, 100663296, -1204394608, 100663296, -1194204647, 100663296, -129639168, 100663296, -1194204645},
			IncidentType:          proxy.IncidentType_MISCELLANEOUS,
			IncidentSeverity:      proxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &proxy.Location{Lat: 44.181220, Lon: -117.135840},
			Description:           "Construction on I-84 EB near MP 359, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "I-84 E",
			EventCode:             500,
			AlertCEventQuantifier: 0,
			IsBlocking:            false,
		},
		&proxy.Incident{
			IncidentId:            "TTI-6f55a1ca-9a6e-38ef-ac40-0dbd3f5586df-TTR83431311705665-1",
			AffectedWayIds:        []int64{100663296, 19446119},
			IncidentType:          proxy.IncidentType_ACCIDENT,
			IncidentSeverity:      proxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &proxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
		},
		&proxy.Incident{
			IncidentId:            "mock-1",
			AffectedWayIds:        []int64{100663296, -1204020275, 100643296},
			IncidentType:          proxy.IncidentType_ACCIDENT,
			IncidentSeverity:      proxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &proxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
		},
	}

	cache := New()

	// update
	cache.Update(newIncidentsResponses(presetIncidents, proxy.Action_UPDATE))
	expectIncidentsCount := 2
	if cache.Count() != expectIncidentsCount {
		t.Errorf("expect cached incidents count %d but got %d", expectIncidentsCount, cache.Count())
	}
	expectAffectedWaysCount := 4 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d", expectAffectedWaysCount, cache.AffectedWaysCount())
	}

	// query expect sucess
	inCacheWayIDs := []int64{100663296, 19446119, -1204020275, 100643296}
	for _, wayID := range inCacheWayIDs {
		if !cache.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect blocked by incident but not", wayID)
		}
	}

	// query expect fail
	notInCacheWayIDs := []int64{0, 100000, -23456789723}
	for _, wayID := range notInCacheWayIDs {
		if cache.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect not blocked by incident but yes", wayID)
		}
	}

	// delete
	deleteIncidents := presetIncidents[:2]
	cache.Update(newIncidentsResponses(deleteIncidents, proxy.Action_DELETE))
	expectIncidentsCount = 1
	if cache.Count() != expectIncidentsCount {
		t.Errorf("expect after delete, cached incidents count %d but got %d", expectIncidentsCount, cache.Count())
	}
	expectAffectedWaysCount = 3 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d", expectAffectedWaysCount, cache.AffectedWaysCount())
	}

	// clear
	cache.Clear()
	if cache.Count() != 0 {
		t.Errorf("expect cached incidents count 0 due to clear but got %d", cache.Count())
	}

}

func newIncidentsResponses(incidents []*proxy.Incident, action proxy.Action) []*proxy.IncidentResponse {

	incidentsResponses := []*proxy.IncidentResponse{}
	for _, incident := range incidents {
		incidentsResponses = append(incidentsResponses, &proxy.IncidentResponse{Incident: incident, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return incidentsResponses
}
