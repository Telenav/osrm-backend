package flowscache

import (
	"reflect"
	"testing"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

func TestFlowsCache(t *testing.T) {

	presetFlows := []*proxy.Flow{
		&proxy.Flow{WayId: -1112859596, Speed: 6.110000, TrafficLevel: proxy.TrafficLevel_SLOW_SPEED},
		&proxy.Flow{WayId: 119961953, Speed: 10.550000, TrafficLevel: proxy.TrafficLevel_SLOW_SPEED},
		&proxy.Flow{WayId: -112614307, Speed: 16.110001, TrafficLevel: proxy.TrafficLevel_FREE_FLOW},
	}

	cache := New()

	// update
	cache.Update(newFlowResponses(presetFlows, proxy.Action_UPDATE))
	if cache.Count() != int64(len(presetFlows)) {
		t.Errorf("expect cached flows count %d but got %d", len(presetFlows), cache.Count())
	}

	// query expect sucess
	for _, f := range presetFlows {
		r := cache.Query(f.WayId)
		if !reflect.DeepEqual(r, f) {
			t.Errorf("Query Flow for wayID %d, expect %v but got %v", f.WayId, f, r)
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

	// delete
	deleteCount := 2
	deleteFlows := presetFlows[:deleteCount]
	cache.Update(newFlowResponses(deleteFlows, proxy.Action_DELETE))
	if cache.Count() != int64(len(presetFlows)-deleteCount) {
		t.Errorf("expect after delete, cached flows count %d but got %d", len(presetFlows)-deleteCount, cache.Count())
	}

	// clear
	cache.Clear()
	if cache.Count() != 0 {
		t.Errorf("expect cached flows count 0 due to clear but got %d", cache.Count())
	}

}

func newFlowResponses(flows []*proxy.Flow, action proxy.Action) []*proxy.FlowResponse {

	flowResponses := []*proxy.FlowResponse{}
	for _, f := range flows {
		flowResponses = append(flowResponses, &proxy.FlowResponse{Flow: f, Action: action, XXX_NoUnkeyedLiteral: struct{}{}, XXX_unrecognized: nil, XXX_sizecache: 0})
	}
	return flowResponses
}
