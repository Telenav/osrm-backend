package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

func TestAddAndGetStartAndEndNodeForNodeGraph(t *testing.T) {
	graph := NewNodeGraph(nil, nil)

	expectStartChargeState := chargingstrategy.State{Energy: 10.0}
	expectStartLocation := &nav.Location{
		Lat: 33.33,
		Lon: -122.22}
	expectEndChargeState := chargingstrategy.State{}
	expectEndLocation := &nav.Location{
		Lat: 34.44,
		Lon: -124.44}

	graph.SetStart(stationfindertype.OrigLocationID.String(), expectStartChargeState, expectStartLocation).
		SetEnd(stationfindertype.DestLocationID.String(), expectEndChargeState, expectEndLocation)

	if graph.StationID(graph.StartNodeID()) != stationfindertype.OrigLocationID.String() {
		t.Errorf("Incorrect result for start's stationID, expect %s but got %s.\n", stationfindertype.OrigLocationID.String(), graph.StationID(graph.StartNodeID()))
	}

	if graph.StationID(graph.EndNodeID()) != stationfindertype.DestLocationID.String() {
		t.Errorf("Incorrect result for end's stationID, expect %s but got %s.\n", stationfindertype.DestLocationID.String(), graph.StationID(graph.EndNodeID()))
	}

	startNode := graph.Node(graph.StartNodeID())
	// if !(reflect.DeepEqual(startNode.Lat, expectStartLocation.Lat) &&
	// 	reflect.DeepEqual(startNode.Lon, expectStartLocation.Lon)) {
	// 	t.Errorf("Incorrect result for start's location, expect %#v but got %#v.\n", expectStartLocation, startNode)
	// }
	if !reflect.DeepEqual(startNode.chargeInfo.targetState, expectStartChargeState) {
		t.Errorf("Incorrect result for start's charge state, expect %#v but got %#v.\n", expectStartChargeState, startNode.chargeInfo.targetState)
	}

	endNode := graph.Node(graph.EndNodeID())
	// if !(reflect.DeepEqual(endNode.Lat, expectEndLocation.Lat) &&
	// 	reflect.DeepEqual(endNode.Lon, expectEndLocation.Lon)) {
	// 	t.Errorf("Incorrect result for end's location, expect %#v but got %#v.\n", expectEndLocation, endNode)
	// }
	if !reflect.DeepEqual(endNode.chargeInfo.targetState, expectEndChargeState) {
		t.Errorf("Incorrect result for end's charge state, expect %#v but got %#v.\n", expectEndChargeState, endNode.chargeInfo.targetState)
	}
}

func TestAdjacentNodesInterfaceForNodeGraph(t *testing.T) {
	graph := generateMockNodeGraph()
	nodeIDs := graph.AdjacentNodes(graph.StartNodeID())

	// testStationID1 connects to 3 physical nodes: test_station_id_2, test_station_id_3, test_station_id_4
	// for each physical node will generate 3 logic nodes: charge to 60%, charge to 80%, charge to full
	if len(nodeIDs) != 9 {
		t.Errorf("Incorrect count of logic nodes generated, expect result is 9 but got %d.\n", len(nodeIDs))
	}

	expectNodes := []*node{
		{
			1,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		{
			2,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		{
			3,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		{
			4,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		{
			5,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		{
			6,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		{
			7,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
		{
			8,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
		{
			9,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
	}

	for i, nodeID := range nodeIDs {
		actualNode := graph.Node(nodeID)
		if !reflect.DeepEqual(actualNode, expectNodes[i]) {
			t.Errorf("Incorrect node generated expect %#v but got %#v.\n", expectNodes[i], actualNode)
		}
	}
}

func TestEdgeInterfaceForNodeGraph(t *testing.T) {
	graph := generateMockNodeGraph()
	nodeIDs := graph.AdjacentNodes(graph.StartNodeID())

	expectEdges := []*edgeMetric{
		{
			distance: 2,
			duration: 2,
		},
		{
			distance: 2,
			duration: 2,
		},
		{
			distance: 2,
			duration: 2,
		},
		{
			distance: 3,
			duration: 3,
		},
		{
			distance: 3,
			duration: 3,
		},
		{
			distance: 3,
			duration: 3,
		},
		{
			distance: 4,
			duration: 4,
		},
		{
			distance: 4,
			duration: 4,
		},
		{
			distance: 4,
			duration: 4,
		},
	}

	actualEdges := make([]*edgeMetric, 0, len(nodeIDs))
	fromNodeID := graph.StartNodeID()
	for _, toNodeID := range nodeIDs {
		actualEdges = append(actualEdges, graph.Edge(fromNodeID, toNodeID))
	}

	if !reflect.DeepEqual(actualEdges, expectEdges) {
		t.Errorf("Incorrect Edge() result for NodeGraph, expect %#v but got %#v\n", expectEdges, actualEdges)
	}
}

func TestStationIDInterfaceForNodeGraph(t *testing.T) {
	graph := generateMockNodeGraph()
	nodeIDs := graph.AdjacentNodes(graph.StartNodeID())

	expectStationIDs := []string{
		testStationID2.String(),
		testStationID2.String(),
		testStationID2.String(),
		testStationID3.String(),
		testStationID3.String(),
		testStationID3.String(),
		testStationID4.String(),
		testStationID4.String(),
		testStationID4.String(),
	}

	actualStationIDs := make([]string, 0, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		actualStationIDs = append(actualStationIDs, graph.StationID(nodeID))
	}
	if !reflect.DeepEqual(actualStationIDs, expectStationIDs) {
		t.Errorf("Incorrect StationID() result for NodeGraph, expect %#v but got %#v\n", expectStationIDs, actualStationIDs)
	}
}

func generateMockNodeGraph() Graph {
	maxEnergyLevel := 100.0
	currEnergyLevel := 10.0
	strategy := chargingstrategy.NewFakeChargingStrategy(maxEnergyLevel)
	querier := newMockQuerier()
	graph := NewNodeGraph(strategy, querier)

	origLocation := querier.GetLocation(testStationID1Str)
	graph.SetStart(testStationID1Str, chargingstrategy.State{Energy: currEnergyLevel}, &nav.Location{Lat: origLocation.Lat, Lon: origLocation.Lon})

	return graph
}

type mockQuerier4NodeGraph struct {
	mockStationID2QueryResult map[string][]*common.RankedPlaceInfo
	mockStationID2Location    map[string]*nav.Location
}

func newMockQuerier() connectivitymap.Querier {
	return &mockQuerier4NodeGraph{
		mockStationID2QueryResult: map[string][]*common.RankedPlaceInfo{
			testStationID1Str: {
				{
					PlaceInfo: common.PlaceInfo{
						ID:       testStationID2,
						Location: &nav.Location{Lat: 2.2, Lon: 2.2},
					},
					Weight: &common.Weight{
						Distance: 2,
						Duration: 2,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       testStationID3,
						Location: &nav.Location{Lat: 3.3, Lon: 3.3},
					},
					Weight: &common.Weight{
						Distance: 3,
						Duration: 3,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       testStationID4,
						Location: &nav.Location{Lat: 4.4, Lon: 4.4},
					},
					Weight: &common.Weight{
						Distance: 4,
						Duration: 4,
					},
				},
			},
		},
		mockStationID2Location: map[string]*nav.Location{
			testStationID1Str: {Lat: 1.1, Lon: 1.1},
			testStationID2Str: {Lat: 2.2, Lon: 2.2},
			testStationID3Str: {Lat: 3.3, Lon: 3.3},
			testStationID4Str: {Lat: 4.4, Lon: 4.4},
		},
	}
}

func (querier *mockQuerier4NodeGraph) NearByStationQuery(stationID string) []*common.RankedPlaceInfo {
	if stationID == testStationID1Str {
		return querier.mockStationID2QueryResult[testStationID1Str]
	}
	glog.Fatalf("Un-implemented mapping key %v for mockStationID2QueryResult.\n", stationID)
	return nil
}

func (querier *mockQuerier4NodeGraph) GetLocation(stationID string) *nav.Location {
	if stationID == testStationID1Str {
		return querier.mockStationID2Location[testStationID1Str]
	}
	glog.Fatalf("Un-implemented mapping key for mockStationID2Location id = %v.\n", stationID)
	return nil
}

var testStationID1Str = stationfindertype.OrigLocationID.String()
var testStationID2Str = "test_station_id_2"
var testStationID3Str = "test_station_id_3"
var testStationID4Str = "test_station_id_4"

var testStationID1 = stationfindertype.OrigLocationID
var testStationID2 common.PlaceID = 2
var testStationID3 common.PlaceID = 3
var testStationID4 common.PlaceID = 4
