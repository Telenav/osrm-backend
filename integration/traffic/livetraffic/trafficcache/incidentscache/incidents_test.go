package incidentscache

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

func TestIncidentsCache(t *testing.T) {
	presetIncidents := []*trafficproxy.Incident{
		&trafficproxy.Incident{
			IncidentID:            "TTI-f47b8dba-59a3-372d-9cec-549eb252e2d5-TTR46312939215361-1",
			AffectedWayIDs:        []int64{100663296, -1204020275, 100663296, -1204020274, 100663296, -916744017, 100663296, -1204020245, 100663296, -1194204646, 100663296, -1204394608, 100663296, -1194204647, 100663296, -129639168, 100663296, -1194204645},
			IncidentType:          trafficproxy.IncidentType_MISCELLANEOUS,
			IncidentSeverity:      trafficproxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &trafficproxy.Location{Lat: 44.181220, Lon: -117.135840},
			Description:           "Construction on I-84 EB near MP 359, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "I-84 E",
			EventCode:             500,
			AlertCEventQuantifier: 0,
			IsBlocking:            false,
			Timestamp:             1579419488000,
		},
		&trafficproxy.Incident{
			IncidentID:            "TTI-6f55a1ca-9a6e-38ef-ac40-0dbd3f5586df-TTR83431311705665-1",
			AffectedWayIDs:        []int64{100663296, 19446119},
			IncidentType:          trafficproxy.IncidentType_ACCIDENT,
			IncidentSeverity:      trafficproxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &trafficproxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
			Timestamp:             1579419488000,
		},
		&trafficproxy.Incident{
			IncidentID:            "mock-1",
			AffectedWayIDs:        []int64{100663296, -1204020275, 100643296},
			IncidentType:          trafficproxy.IncidentType_ACCIDENT,
			IncidentSeverity:      trafficproxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &trafficproxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
			Timestamp:             1579419488000,
		},
	}

	updateIncidents := []*trafficproxy.Incident{
		&trafficproxy.Incident{
			IncidentID:            "mock-1",
			AffectedWayIDs:        []int64{100663296, -1204020275, 100643296, 111111111},
			IncidentType:          trafficproxy.IncidentType_ACCIDENT,
			IncidentSeverity:      trafficproxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &trafficproxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
			Timestamp:             1579419500000,
		}, // newer
		&trafficproxy.Incident{
			IncidentID:            "mock-1",
			AffectedWayIDs:        []int64{100663296, -1204020275, 100643296},
			IncidentType:          trafficproxy.IncidentType_ACCIDENT,
			IncidentSeverity:      trafficproxy.IncidentSeverity_CRITICAL,
			IncidentLocation:      &trafficproxy.Location{Lat: 37.592370, Lon: -77.56735040},
			Description:           "Incident on N PARHAM RD near RIDGE RD, Drive with caution.",
			FirstCrossStreet:      "",
			SecondCrossStreet:     "",
			StreetName:            "N Parham Rd",
			EventCode:             214,
			AlertCEventQuantifier: 0,
			IsBlocking:            true,
			Timestamp:             1579419000000,
		}, // older
	}

	cache := New()

	// update
	cache.Update(newIncidentsResponses(presetIncidents, trafficproxy.Action_UPDATE))
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

	// update
	cache.Update(newIncidentsResponses(updateIncidents, trafficproxy.Action_UPDATE))
	if cache.Count() != expectIncidentsCount {
		t.Errorf("expect cached incidents count %d but got %d", expectIncidentsCount, cache.Count())
	}
	expectAffectedWaysCount = 5 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d", expectAffectedWaysCount, cache.AffectedWaysCount())
	}

	// query expect sucess
	inCacheWayIDs = []int64{111111111} // only check the updated one
	for _, wayID := range inCacheWayIDs {
		if !cache.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect blocked by incident but not", wayID)
		}
	}

	// delete
	deleteIncidents := presetIncidents[:2]
	cache.Update(newIncidentsResponses(deleteIncidents, trafficproxy.Action_DELETE))
	expectIncidentsCount = 1
	if cache.Count() != expectIncidentsCount {
		t.Errorf("expect after delete, cached incidents count %d but got %d", expectIncidentsCount, cache.Count())
	}
	expectAffectedWaysCount = 4 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d", expectAffectedWaysCount, cache.AffectedWaysCount())
	}

	// clear
	cache.Clear()
	if cache.Count() != 0 {
		t.Errorf("expect cached incidents count 0 due to clear but got %d", cache.Count())
	}
	if cache.AffectedWaysCount() != 0 {
		t.Errorf("expect cached incidents affect ways count 0 but got %d", cache.AffectedWaysCount())
	}

}

func newIncidentsResponses(incidents []*trafficproxy.Incident, action trafficproxy.Action) []*trafficproxy.IncidentResponse {

	incidentsResponses := []*trafficproxy.IncidentResponse{}
	for _, incident := range incidents {
		incidentsResponses = append(incidentsResponses, &trafficproxy.IncidentResponse{Incident: incident, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return incidentsResponses
}
