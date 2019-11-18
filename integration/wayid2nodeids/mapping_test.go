package wayid2nodeids

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/nodebasededge"
)

func TestParseLine(t *testing.T) {

	cases := []struct {
		line string
		ids  []int64
	}{
		{"", nil},
		{"24418325,84760891102", nil},
		{"24418325,84760891102,19496208102", []int64{24418325, 84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,", []int64{24418325, 84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,,,,,", []int64{24418325, 84760891102, 19496208102}},
	}

	for _, c := range cases {
		result := parseLine(c.line)
		if !reflect.DeepEqual(result, c.ids) {
			t.Errorf("parseLine %s, expect %v, but got %v", c.line, c.ids, result)
		}
	}
}

func TestMappingLoad(t *testing.T) {

	m := NewMappingFrom("./testdata/sample_wayid2nodeids.csv.snappy")
	if err := m.Load(); err != nil {
		t.Error(err)
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

	expectEdge2WayIDMapping := map[nodebasededge.Edge]int64{
		nodebasededge.Edge{FromNode: 84760891102, ToNode: 19496208102}:     24418325,
		nodebasededge.Edge{FromNode: 84762609102, ToNode: 244183320001101}: 24418332,
		nodebasededge.Edge{FromNode: 244183320001101, ToNode: 84762607102}: 24418332,
		nodebasededge.Edge{FromNode: 84760849102, ToNode: 84760850102}:     24418343,
		nodebasededge.Edge{FromNode: 84760846102, ToNode: 84760858102}:     24418344,
	}

	if !reflect.DeepEqual(expectEdge2WayIDMapping, m.edge2WayID) {
		t.Errorf("expect edge2wayID mapping %v, but got %v", expectEdge2WayIDMapping, m.edge2WayID)
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
		gotNodeIDs := m.GetNodeIDs(c.wayID)
		if !reflect.DeepEqual(gotNodeIDs, c.expectNodeIDs) {
			t.Errorf("expect nodeIDs %v for wayID %d, but got %v", c.expectNodeIDs, c.wayID, gotNodeIDs)
		}
	}

	// GetWayID
	getWayIDCases := []struct {
		edge        nodebasededge.Edge
		expectGot   bool
		expectWayID int64
	}{
		{nodebasededge.Edge{FromNode: 12345, ToNode: 67890}, false, 0},
		{nodebasededge.Edge{FromNode: 84762609102, ToNode: 244183320001101}, true, 24418332},
		{nodebasededge.Edge{FromNode: 244183320001101, ToNode: 84762607102}, true, 24418332},
		{nodebasededge.Edge{FromNode: 84762607102, ToNode: 244183320001101}, true, -24418332},
	}

	for _, c := range getWayIDCases {
		gotWayID, ok := m.GetWayID(c.edge)
		if c.expectGot != ok || (ok && gotWayID != c.expectWayID) {
			t.Errorf("expect wayID %d,%t for Edge %v, but got %d,%t", c.expectWayID, c.expectGot, c.edge, gotWayID, ok)
		}
	}

}
