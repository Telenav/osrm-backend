package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

func TestAddAndGetStartAndEndNodeForNodeGraph(t *testing.T) {
	graph := NewNodeGraph(nil, nil)

	expectStartChargeState := chargingstrategy.State{Energy: 10.0}
	expectStartLocation := locationInfo{33.33, -122.22}
	expectEndChargeState := chargingstrategy.State{}
	expectEndLocation := locationInfo{34.44, -124.44}

	graph.SetStart(stationfindertype.OrigLocationID, expectStartChargeState, expectStartLocation).
		SetEnd(stationfindertype.DestLocationID, expectEndChargeState, expectEndLocation)

	if graph.StationID(graph.StartNodeID()) != stationfindertype.OrigLocationID {
		t.Errorf("Incorrect result for start's stationID, expect %s but got %s.\n", stationfindertype.OrigLocationID, graph.StationID(graph.StartNodeID()))
	}

	if graph.StationID(graph.EndNodeID()) != stationfindertype.DestLocationID {
		t.Errorf("Incorrect result for end's stationID, expect %s but got %s.\n", stationfindertype.DestLocationID, graph.StationID(graph.EndNodeID()))
	}

	startNode := graph.Node(graph.StartNodeID())
	if !reflect.DeepEqual(startNode.locationInfo, expectStartLocation) {
		t.Errorf("Incorrect result for start's location, expect %#v but got %#v.\n", expectStartLocation, startNode.locationInfo)
	}
	if !reflect.DeepEqual(startNode.chargeInfo.targetState, expectStartChargeState) {
		t.Errorf("Incorrect result for start's charge state, expect %#v but got %#v.\n", expectStartChargeState, startNode.chargeInfo.targetState)
	}

	endNode := graph.Node(graph.EndNodeID())
	if !reflect.DeepEqual(endNode.locationInfo, expectEndLocation) {
		t.Errorf("Incorrect result for end's location, expect %#v but got %#v.\n", expectEndLocation, endNode.locationInfo)
	}
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
			id: 1,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			locationInfo: locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			id: 2,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			locationInfo: locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			id: 3,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			locationInfo: locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			id: 4,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			locationInfo: locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			id: 5,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			locationInfo: locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			id: 6,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			locationInfo: locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			id: 7,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 60.0,
				},
			},
			locationInfo: locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
		},
		{
			id: 8,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 80.0,
				},
			},
			locationInfo: locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
		},
		{
			id: 9,
			chargeInfo: chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 100.0,
				},
			},
			locationInfo: locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
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

	expectEdges := []*edge{
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

	actualEdges := make([]*edge, 0, len(nodeIDs))
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

	origLocation := querier.GetLocation(testStationID1)
	graph.SetStart(testStationID1, chargingstrategy.State{Energy: currEnergyLevel}, locationInfo{lat: origLocation.Lat, lon: origLocation.Lon})

	return graph
}

type mockQuerier4NodeGraph struct {
	mockStationID2QueryResult map[string][]*connectivitymap.QueryResult
	mockStationID2Location    map[string]*nav.Location
}

func newMockQuerier() connectivitymap.Querier {
	return &mockQuerier4NodeGraph{
		mockStationID2QueryResult: map[string][]*connectivitymap.QueryResult{
			testStationID1: {
				{
					StationID:       testStationID2,
					StationLocation: &nav.Location{Lat: 2.2, Lon: 2.2},
					Distance:        2,
					Duration:        2,
				},
				{
					StationID:       testStationID3,
					StationLocation: &nav.Location{Lat: 3.3, Lon: 3.3},
					Distance:        3,
					Duration:        3,
				},
				{
					StationID:       testStationID4,
					StationLocation: &nav.Location{Lat: 4.4, Lon: 4.4},
					Distance:        4,
					Duration:        4,
				},
			},
		},
		mockStationID2Location: map[string]*nav.Location{
			testStationID1: {Lat: 1.1, Lon: 1.1},
			testStationID2: {Lat: 2.2, Lon: 2.2},
			testStationID3: {Lat: 3.3, Lon: 3.3},
			testStationID4: {Lat: 4.4, Lon: 4.4},
		},
	}
}

func (querier *mockQuerier4NodeGraph) NearByStationQuery(stationID string) []*connectivitymap.QueryResult {
	if stationID == testStationID1 {
		return querier.mockStationID2QueryResult[testStationID1]
	}
	glog.Fatal("Un-implemented mapping key for mockStationID2QueryResult.\n")
	return nil
}

func (querier *mockQuerier4NodeGraph) GetLocation(stationID string) *nav.Location {
	if stationID == testStationID1 {
		return querier.mockStationID2Location[testStationID1]
	}
	glog.Fatal("Un-implemented mapping key for mockStationID2Location.\n")
	return nil
}

var testStationID1 = stationfindertype.OrigLocationID
var testStationID2 = "test_station_id_2"
var testStationID3 = "test_station_id_3"
var testStationID4 = "test_station_id_4"
