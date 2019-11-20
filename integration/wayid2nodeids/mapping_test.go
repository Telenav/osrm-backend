package wayid2nodeids

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/graph"
)

func TestMappingLoad(t *testing.T) {

	m := NewMappingFrom("./testdata/sample_wayid2nodeids.csv.snappy")
	if err := m.Load(); err != nil {
		t.Error(err)
	}

	if !m.IsReady() {
		t.Error("expect ready but not")
	}

	expectWayID2NodeIDsMapping := map[int64][]int64{
		24418325: []int64{84760891102, 19496208102},
		24418332: []int64{84762609102, 244183320001101, 84762607102},
		24418343: []int64{84760849102, 84760850102},
		24418344: []int64{84760846102, 84760858102},
	}

	if !reflect.DeepEqual(expectWayID2NodeIDsMapping, m.wayID2NodeIDs) {
		t.Errorf("expect wayid2nodeids mapping %v, but got %v", expectWayID2NodeIDsMapping, m.wayID2NodeIDs)
	}

	// GetNodeIDs
	getNodesCases := []struct {
		wayID         int64
		expectNodeIDs []int64
	}{
		{240000, nil},
		{24418325, []int64{84760891102, 19496208102}},
		{24418332, []int64{84762609102, 244183320001101, 84762607102}},
	}
	for _, c := range getNodesCases {
		gotNodeIDs := m.WayID2NodeIDs(c.wayID)
		if !reflect.DeepEqual(gotNodeIDs, c.expectNodeIDs) {
			t.Errorf("expect nodeIDs %v for wayID %d, but got %v", c.expectNodeIDs, c.wayID, gotNodeIDs)
		}
	}

	// GetEdges
	getEdgesCases := []struct {
		wayID       int64
		expectEdges []graph.Edge
	}{
		{240000, nil},
		{24418325, []graph.Edge{graph.Edge{From: 84760891102, To: 19496208102}}},
		{
			24418332, []graph.Edge{
				graph.Edge{From: 84762609102, To: 244183320001101},
				graph.Edge{From: 244183320001101, To: 84762607102}},
		},
		{
			-24418332, []graph.Edge{
				graph.Edge{From: 84762607102, To: 244183320001101},
				graph.Edge{From: 244183320001101, To: 84762609102}},
		},
	}
	for _, c := range getEdgesCases {
		gotEdges := m.WayID2Edges(c.wayID)
		if !reflect.DeepEqual(gotEdges, c.expectEdges) {
			t.Errorf("expect edges %v for wayID %d, but got %v", c.expectEdges, c.wayID, gotEdges)
		}
	}

}
