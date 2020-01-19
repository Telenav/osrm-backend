package flowscache

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
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

	cache := New()

	// update preset
	cache.Update(newFlowResponses(presetFlows, trafficproxy.Action_UPDATE))
	if cache.Count() != int64(len(presetFlows)) {
		t.Errorf("expect cached flows count %d but got %d", len(presetFlows), cache.Count())
	}

	// query expect sucess
	for _, f := range presetFlows {
		r := cache.Query(f.WayID)
		if !reflect.DeepEqual(r, f) {
			t.Errorf("Query Flow for wayID %d, expect %v but got %v", f.WayID, f, r)
		}
	}

	// query expect fail
	notInCacheWayIDs := []int64{0, 100000, -23456789723}
	for _, wayID := range notInCacheWayIDs {
		r := cache.Query(wayID)
		if r != nil {
			t.Errorf("Query Flow for wayID %d, expect nil but got %v", wayID, r)
		}
	}

	// update exists
	cache.Update(newFlowResponses(updateFlows, trafficproxy.Action_UPDATE))
	if cache.Count() != int64(len(presetFlows)) { // expect no change
		t.Errorf("expect cached flows count %d but got %d", len(presetFlows), cache.Count())
	}

	// query expect sucess
	queryExpectFlow := updateFlows[0]
	r := cache.Query(queryExpectFlow.WayID)
	if !reflect.DeepEqual(r, queryExpectFlow) {
		t.Errorf("Query Flow for wayID %d, expect %v but got %v", queryExpectFlow.WayID, queryExpectFlow, r)
	}

	// delete
	deleteCount := 2
	deleteFlows := presetFlows[:deleteCount]
	cache.Update(newFlowResponses(deleteFlows, trafficproxy.Action_DELETE))
	if cache.Count() != int64(len(presetFlows)-deleteCount) {
		t.Errorf("expect after delete, cached flows count %d but got %d", len(presetFlows)-deleteCount, cache.Count())
	}

	// clear
	cache.Clear()
	if cache.Count() != 0 {
		t.Errorf("expect cached flows count 0 due to clear but got %d", cache.Count())
	}

}

func newFlowResponses(flows []*trafficproxy.Flow, action trafficproxy.Action) []*trafficproxy.FlowResponse {

	flowResponses := []*trafficproxy.FlowResponse{}
	for _, f := range flows {
		flowResponses = append(flowResponses, &trafficproxy.FlowResponse{Flow: f, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return flowResponses
}
