package flowscacheindexedbyedge

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/graph"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

func TestFlowsCache(t *testing.T) {

	presetFlows := []*trafficproxy.Flow{
		&trafficproxy.Flow{WayID: -1112859596, Speed: 6.110000, TrafficLevel: trafficproxy.TrafficLevel_SLOW_SPEED, Timestamp: 1579419488000},
		&trafficproxy.Flow{WayID: 119961953, Speed: 10.550000, TrafficLevel: trafficproxy.TrafficLevel_SLOW_SPEED, Timestamp: 1579419488000},
		&trafficproxy.Flow{WayID: -112614307, Speed: 16.110001, TrafficLevel: trafficproxy.TrafficLevel_FREE_FLOW, Timestamp: 1579419488000},
	}

	updateFlows := []*trafficproxy.Flow{
		&trafficproxy.Flow{WayID: -112614307, Speed: 20.110001, TrafficLevel: trafficproxy.TrafficLevel_FREE_FLOW, Timestamp: 1579419500000}, // newer
		&trafficproxy.Flow{WayID: -112614307, Speed: 13.110001, TrafficLevel: trafficproxy.TrafficLevel_FREE_FLOW, Timestamp: 1579419000000}, // older
	}

	wayid2NodeIDsMapping := wayID2NodeIDs{
		1112859596: []int64{123456789001, 123456789002, 123456789003, 123456789004},
		119961953:  []int64{123456789011, 123456789012, 123456789003, 123456789014, 123456789015},
		112614307:  []int64{123456789021, 123456789002, 123456789023, 123456789024, 123456789025, 123456789026},
	}

	cache := New(wayid2NodeIDsMapping)

	// update
	cache.Update(newFlowResponses(presetFlows, trafficproxy.Action_UPDATE))
	expectedFlowsIndexedByEdgeCount := int64(12)
	if cache.Count() != expectedFlowsIndexedByEdgeCount {
		t.Errorf("expect flows count %d but got %d", expectedFlowsIndexedByEdgeCount, cache.Count())
	}
	if cache.AffectedWaysCount() != int64(len(presetFlows)) {
		t.Errorf("expect flows affected ways count %d but got %d", len(presetFlows), cache.AffectedWaysCount())
	}

	// query expect sucess
	inCacheEdges := []graph.Edge{
		graph.Edge{From: 123456789004, To: 123456789003},
		graph.Edge{From: 123456789011, To: 123456789012},
		graph.Edge{From: 123456789023, To: 123456789002},
	}
	for i := range inCacheEdges {
		r := cache.QueryByEdge(inCacheEdges[i])
		if !reflect.DeepEqual(r, presetFlows[i]) {
			t.Errorf("Query Flow for Edge %v, expect %v but got %v", inCacheEdges[i], presetFlows[i], r)
		}
	}
	gotFlows := cache.QueryByEdges(inCacheEdges)
	if len(gotFlows) != len(inCacheEdges) {
		t.Errorf("Query Flow on Edges %v, expect flows count %d but got %d", inCacheEdges, len(inCacheEdges), len(gotFlows))
	}

	// query expect fail
	notInCacheEdges := []graph.Edge{
		graph.Edge{},
		graph.Edge{From: 12345, To: 123456789004},
		graph.Edge{From: 1000000000, To: 123456789012},
		graph.Edge{From: 123456789001, To: 123456789002},
	}
	for _, e := range notInCacheEdges {
		r := cache.QueryByEdge(e)
		if r != nil {
			t.Errorf("Query Flow on Edge %v, expect nil but got %v", e, r)
		}
	}

	// update exists
	cache.Update(newFlowResponses(updateFlows, trafficproxy.Action_UPDATE))
	if cache.Count() != expectedFlowsIndexedByEdgeCount {
		t.Errorf("expect flows count %d but got %d", expectedFlowsIndexedByEdgeCount, cache.Count())
	}
	if cache.AffectedWaysCount() != int64(len(presetFlows)) { // expect no change
		t.Errorf("expect flows affected ways count %d but got %d", len(presetFlows), cache.AffectedWaysCount())
	}

	// query expect sucess
	queryExpectFlow := updateFlows[0]
	queryEdge := graph.Edge{From: 123456789002, To: 123456789021}
	r := cache.QueryByEdge(queryEdge)
	if !reflect.DeepEqual(r, queryExpectFlow) {
		t.Errorf("Query Flow for Edge %v, expect %v but got %v", queryEdge, queryExpectFlow, r)
	}

	// delete
	deleteCount := 2
	deleteFlows := presetFlows[:deleteCount]
	cache.Update(newFlowResponses(deleteFlows, trafficproxy.Action_DELETE))
	expectedFlowsIndexedByEdgeCount = int64(5)
	if cache.Count() != expectedFlowsIndexedByEdgeCount {
		t.Errorf("expect flows count %d but got %d", expectedFlowsIndexedByEdgeCount, cache.Count())
	}
	if cache.AffectedWaysCount() != int64(len(presetFlows)-deleteCount) {
		t.Errorf("expect after delete, flows affected ways count %d but got %d", len(presetFlows)-deleteCount, cache.AffectedWaysCount())
	}

	// clear
	cache.Clear()
	if cache.Count() != 0 || cache.AffectedWaysCount() != 0 {
		t.Errorf("expect cached flows,affectedways count 0,0 due to clear but got %d,%d", cache.Count(), cache.AffectedWaysCount())
	}

}

func newFlowResponses(flows []*trafficproxy.Flow, action trafficproxy.Action) []*trafficproxy.FlowResponse {

	flowResponses := []*trafficproxy.FlowResponse{}
	for _, f := range flows {
		flowResponses = append(flowResponses, &trafficproxy.FlowResponse{Flow: f, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return flowResponses
}
