package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
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

	graph.SetStart(iteratortype.OrigLocationID, expectStartChargeState, expectStartLocation).
		SetEnd(iteratortype.DestLocationID, expectEndChargeState, expectEndLocation)

	if graph.PlaceID(graph.StartNodeID()) != iteratortype.OrigLocationID {
		t.Errorf("Incorrect result for start's placeID, expect %s but got %s.\n", iteratortype.OrigLocationID.String(), graph.PlaceID(graph.StartNodeID()))
	}

	if graph.PlaceID(graph.EndNodeID()) != iteratortype.DestLocationID {
		t.Errorf("Incorrect result for end's placeID, expect %s but got %s.\n", iteratortype.DestLocationID.String(), graph.PlaceID(graph.EndNodeID()))
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

	expectEdges := []*entity.Weight{
		{
			Distance: 2,
			Duration: 2,
		},
		{
			Distance: 2,
			Duration: 2,
		},
		{
			Distance: 2,
			Duration: 2,
		},
		{
			Distance: 3,
			Duration: 3,
		},
		{
			Distance: 3,
			Duration: 3,
		},
		{
			Distance: 3,
			Duration: 3,
		},
		{
			Distance: 4,
			Duration: 4,
		},
		{
			Distance: 4,
			Duration: 4,
		},
		{
			Distance: 4,
			Duration: 4,
		},
	}

	actualEdges := make([]*entity.Weight, 0, len(nodeIDs))
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

	expectStationIDs := []entity.PlaceID{
		testStationID2,
		testStationID2,
		testStationID2,
		testStationID3,
		testStationID3,
		testStationID3,
		testStationID4,
		testStationID4,
		testStationID4,
	}

	actualStationIDs := make([]entity.PlaceID, 0, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		actualStationIDs = append(actualStationIDs, graph.PlaceID(nodeID))
	}
	if !reflect.DeepEqual(actualStationIDs, expectStationIDs) {
		t.Errorf("Incorrect PlaceID() result for NodeGraph, expect %#v but got %#v\n", expectStationIDs, actualStationIDs)
	}
}

func generateMockNodeGraph() Graph {
	maxEnergyLevel := 100.0
	currEnergyLevel := 10.0
	strategy := chargingstrategy.NewSimpleChargingStrategy(maxEnergyLevel)
	querier := newMockQuerier()
	graph := NewNodeGraph(strategy, querier)

	origLocation := querier.GetLocation(testStationID1)
	graph.SetStart(testStationID1, chargingstrategy.State{Energy: currEnergyLevel}, &nav.Location{Lat: origLocation.Lat, Lon: origLocation.Lon})

	return graph
}

type mockQuerier4NodeGraph struct {
	mockStationID2QueryResult map[entity.PlaceID][]*entity.TransferInfo
	mockStationID2Location    map[entity.PlaceID]*nav.Location
}

func newMockQuerier() place.TopoQuerier {
	return &mockQuerier4NodeGraph{
		mockStationID2QueryResult: map[entity.PlaceID][]*entity.TransferInfo{
			testStationID1: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testStationID2,
						Location: &nav.Location{Lat: 2.2, Lon: 2.2},
					},
					Weight: &entity.Weight{
						Distance: 2,
						Duration: 2,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testStationID3,
						Location: &nav.Location{Lat: 3.3, Lon: 3.3},
					},
					Weight: &entity.Weight{
						Distance: 3,
						Duration: 3,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testStationID4,
						Location: &nav.Location{Lat: 4.4, Lon: 4.4},
					},
					Weight: &entity.Weight{
						Distance: 4,
						Duration: 4,
					},
				},
			},
		},
		mockStationID2Location: map[entity.PlaceID]*nav.Location{
			testStationID1: {Lat: 1.1, Lon: 1.1},
			testStationID2: {Lat: 2.2, Lon: 2.2},
			testStationID3: {Lat: 3.3, Lon: 3.3},
			testStationID4: {Lat: 4.4, Lon: 4.4},
		},
	}
}

func (querier *mockQuerier4NodeGraph) GetConnectedPlaces(placeID entity.PlaceID) []*entity.TransferInfo {
	if placeID == testStationID1 {
		return querier.mockStationID2QueryResult[testStationID1]
	}
	glog.Fatalf("Un-implemented mapping key %v for mockStationID2QueryResult.\n", placeID)
	return nil
}

func (querier *mockQuerier4NodeGraph) GetLocation(placeID entity.PlaceID) *nav.Location {
	if placeID == testStationID1 {
		return querier.mockStationID2Location[testStationID1]
	}
	glog.Fatalf("Un-implemented mapping key for mockStationID2Location id = %v.\n", placeID)
	return nil
}

var testStationID1Str = iteratortype.OrigLocationID.String()
var testStationID2Str = "test_station_id_2"
var testStationID3Str = "test_station_id_3"
var testStationID4Str = "test_station_id_4"

var testStationID1 = iteratortype.OrigLocationID
var testStationID2 entity.PlaceID = 2
var testStationID3 entity.PlaceID = 3
var testStationID4 entity.PlaceID = 4
