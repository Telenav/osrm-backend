package incidentscache

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/graph"
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

	wayid2NodeIDsMapping := wayID2NodeIDs{
		1204020274: []int64{123456789001, 123456789002, 123456789003, 123456789004},
		100663296:  []int64{123456789011, 123456789012, 123456789003},
		19446119:   []int64{123456789021, 123456789002, 123456789023, 123456789024, 123456789025, 123456789026},
		1204020275: []int64{123456789031, 123456789032, 123456789033, 123456789034},
		100643296:  []int64{123456789041, 123456789042, 123456789043},
		111111111:  []int64{11111111101, 1111111102, 1111111103},
	}

	cache := New()
	cacheWithEdgeIndexing := NewWithEdgeIndexing(wayid2NodeIDsMapping)

	// update
	cache.Update(newIncidentsResponses(presetIncidents, trafficproxy.Action_UPDATE))
	cacheWithEdgeIndexing.Update(newIncidentsResponses(presetIncidents, trafficproxy.Action_UPDATE))
	expectIncidentsCount := 2
	if cache.Count() != expectIncidentsCount || cacheWithEdgeIndexing.Count() != expectIncidentsCount {
		t.Errorf("expect cached incidents count %d but got %d,%d", expectIncidentsCount, cache.Count(), cacheWithEdgeIndexing.Count())
	}
	expectAffectedWaysCount := 4 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount || cacheWithEdgeIndexing.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d,%d", expectAffectedWaysCount, cache.AffectedWaysCount(), cacheWithEdgeIndexing.AffectedWaysCount())
	}
	expectAffectedEdgesCount := 12
	if cacheWithEdgeIndexing.AffectedEdgesCount() != expectAffectedEdgesCount {
		t.Errorf("expect cached incidents affect edges count %d but got %d", expectAffectedEdgesCount, cacheWithEdgeIndexing.AffectedEdgesCount())
	}

	// query expect sucess
	inCacheWayIDs := []int64{100663296, 19446119, -1204020275, 100643296}
	for _, wayID := range inCacheWayIDs {
		if !cache.WayBlockedByIncident(wayID) || !cacheWithEdgeIndexing.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect blocked by incident but not", wayID)
		}
		edges := wayid2NodeIDsMapping.WayID2Edges(wayID)
		for _, e := range edges {
			if !cacheWithEdgeIndexing.EdgeBlockedByIncident(e) {
				t.Errorf("edge %v, expect blocked by incident but not", e)
			}
		}
		if b, i := cacheWithEdgeIndexing.EdgesBlockedByIncidents(edges); !b || i != 0 {
			t.Errorf("edges %v, expect blocked by incidents but not", edges)
		}
	}

	// query expect fail
	notInCacheWayIDs := []int64{0, 100000, -23456789723}
	for _, wayID := range notInCacheWayIDs {
		if cache.WayBlockedByIncident(wayID) || cacheWithEdgeIndexing.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect not blocked by incident but yes", wayID)
		}
	}
	notInCacheEdges := []graph.Edge{
		graph.Edge{},
		graph.Edge{From: 12345, To: 123456789004},
		graph.Edge{From: 1000000000, To: 123456789012},
		graph.Edge{From: 123456789001, To: 123456789002},
	}
	for _, e := range notInCacheEdges {
		if cacheWithEdgeIndexing.EdgeBlockedByIncident(e) {
			t.Errorf("edge %v, expect not blocked by incident but yes", e)
		}
	}

	// update
	cache.Update(newIncidentsResponses(updateIncidents, trafficproxy.Action_UPDATE))
	cacheWithEdgeIndexing.Update(newIncidentsResponses(updateIncidents, trafficproxy.Action_UPDATE))
	if cache.Count() != expectIncidentsCount || cacheWithEdgeIndexing.Count() != expectIncidentsCount {
		t.Errorf("expect cached incidents count %d but got %d,%d", expectIncidentsCount, cache.Count(), cacheWithEdgeIndexing.Count())
	}
	expectAffectedWaysCount = 5 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount || cacheWithEdgeIndexing.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d,%d", expectAffectedWaysCount, cache.AffectedWaysCount(), cacheWithEdgeIndexing.AffectedWaysCount())
	}
	expectAffectedEdgesCount = 14
	if cacheWithEdgeIndexing.AffectedEdgesCount() != expectAffectedEdgesCount {
		t.Errorf("expect cached incidents affect edges count %d but got %d", expectAffectedEdgesCount, cacheWithEdgeIndexing.AffectedEdgesCount())
	}

	// query expect sucess
	inCacheWayIDs = []int64{111111111} // only check the updated one
	for _, wayID := range inCacheWayIDs {
		if !cache.WayBlockedByIncident(wayID) || !cacheWithEdgeIndexing.WayBlockedByIncident(wayID) {
			t.Errorf("wayID %d, expect blocked by incident but not", wayID)
		}
		edges := wayid2NodeIDsMapping.WayID2Edges(wayID)
		for _, e := range edges {
			if !cacheWithEdgeIndexing.EdgeBlockedByIncident(e) {
				t.Errorf("edge %v, expect blocked by incident but not", e)
			}
		}
		if b, i := cacheWithEdgeIndexing.EdgesBlockedByIncidents(edges); !b || i != 0 {
			t.Errorf("edges %v, expect blocked by incidents but not", edges)
		}
	}

	// delete
	deleteIncidents := presetIncidents[:2]
	cache.Update(newIncidentsResponses(deleteIncidents, trafficproxy.Action_DELETE))
	cacheWithEdgeIndexing.Update(newIncidentsResponses(deleteIncidents, trafficproxy.Action_DELETE))
	expectIncidentsCount = 1
	if cache.Count() != expectIncidentsCount || cacheWithEdgeIndexing.Count() != expectIncidentsCount {
		t.Errorf("expect after delete, cached incidents count %d but got %d,%d", expectIncidentsCount, cache.Count(), cacheWithEdgeIndexing.Count())
	}
	expectAffectedWaysCount = 4 // only store blocked incidents
	if cache.AffectedWaysCount() != expectAffectedWaysCount || cacheWithEdgeIndexing.AffectedWaysCount() != expectAffectedWaysCount {
		t.Errorf("expect cached incidents affect ways count %d but got %d,%d", expectAffectedWaysCount, cache.AffectedWaysCount(), cacheWithEdgeIndexing.AffectedWaysCount())
	}
	expectAffectedEdgesCount = 9
	if cacheWithEdgeIndexing.AffectedEdgesCount() != expectAffectedEdgesCount {
		t.Errorf("expect cached incidents affect edges count %d but got %d", expectAffectedEdgesCount, cacheWithEdgeIndexing.AffectedEdgesCount())
	}

	// clear
	cache.Clear()
	cacheWithEdgeIndexing.Clear()
	if cache.Count() != 0 || cacheWithEdgeIndexing.Count() != 0 {
		t.Errorf("expect cached incidents count 0 due to clear but got %d,%d", cache.Count(), cacheWithEdgeIndexing.Count())
	}
	if cache.AffectedWaysCount() != 0 || cacheWithEdgeIndexing.AffectedWaysCount() != 0 {
		t.Errorf("expect cached incidents affect ways count 0 but got %d,%d", cache.AffectedWaysCount(), cacheWithEdgeIndexing.AffectedWaysCount())
	}
	if cacheWithEdgeIndexing.AffectedEdgesCount() != 0 {
		t.Errorf("expect cached incidents affect edges count 0 but got %d", cacheWithEdgeIndexing.AffectedEdgesCount())
	}

}

func newIncidentsResponses(incidents []*trafficproxy.Incident, action trafficproxy.Action) []*trafficproxy.IncidentResponse {

	incidentsResponses := []*trafficproxy.IncidentResponse{}
	for _, incident := range incidents {
		incidentsResponses = append(incidentsResponses, &trafficproxy.IncidentResponse{Incident: incident, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return incidentsResponses
}
