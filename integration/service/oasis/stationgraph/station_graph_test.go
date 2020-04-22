package stationgraph

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

// import (
// 	"fmt"
// 	"math"
// 	"reflect"
// 	"testing"

// 	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
// 	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
// 	"github.com/Telenav/osrm-backend/integration/service/oasis/solution"
// 	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
// 	"github.com/Telenav/osrm-backend/integration/util"
// )

// /*
// Construct test graph of
// - start connects to staion 1, station 2, station 3
//      - start -> station 1: 22.2,
//      - start -> station 2: 11.1,
//      - start -> station 3: 33.3,
// - station 1 connects to station 4, station 5,
//      - station 1 -> station 4: 44.4,
//      - station 1 -> station 5: 34.4,
// - station 2 connects to station 4, station 5
//      - station 2 -> station 4: 11.1,
//      - station 2 -> station 5: 14.4,
// - station 3 connects to station 4, station 5
//      - station 3 -> station 4: 22.2,
//      - station 3 -> station 5: 15.5,
// - station 4 connects to end
//      - station 4 -> end      : 44.4,
// - station 5 connects to end
//      - station 5 -> end      : 33.3,
//                     station 1
//                /       \    \
//               /         \    \
//              /          _\_____station 4
//             /          /  \     /        \
//            /          /    \   /          \
// start  -------   station 2  \ /           end
//            \          \     / \          /
//             \          \   /   \        /
//              \          \_/_____station 5
//               \          /     /
//                \        /     /
//                     station 3

// - Each charge station will try to generate up to 3 virtual node based on fakechargestrategy,
//   each node represent for situation of time spend in charging and get to different energy level
// - In total there could be 17 nodes in graph
// 	  + start node
// 	  + end node
// 	  + station 1 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
// 	  + station 2 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
// 	  + station 3 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
// 	  + station 4 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
// 	  + station 5 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
// - Take station 1 as an example, it will create 3 nodes. For each of them, will contains neighbor information with:
// 	  + connection to station 4's node with 60% of total energy
// 	  + connection to station 4's node with 80% of total energy
// 	  + connection to station 4's node with 100% of total energy
// 	  + connection to station 5's node with 60% of total energy
// 	  + connection to station 5's node with 80% of total energy
// 	  + connection to station 5's node with 100% of total energy
//   Each node with name `station1` will have different id, and each of them will have 6 neighbor nodes
// - For default graph, initial energy is 20.0, max energy is 50.0
// */

var testSGStationID1 = "station1"
var testSGStationID2 = "station2"
var testSGStationID3 = "station3"
var testSGStationID4 = "station4"
var testSGStationID5 = "station5"

type mockQuerier4StationGraph struct {
	mockStationID2QueryResult map[string][]*connectivitymap.QueryResult
	mockStationID2Location    map[string]*nav.Location
}

func newMockQuerier4StationGraph() connectivitymap.Querier {
	querier := &mockQuerier4StationGraph{
		mockStationID2QueryResult: map[string][]*connectivitymap.QueryResult{
			stationfindertype.OrigLocationID: {
				{
					StationID:       testSGStationID1,
					StationLocation: &nav.Location{Lat: 1.1, Lon: 1.1},
					Distance:        22.2,
					Duration:        22.2,
				},
				{
					StationID:       testSGStationID2,
					StationLocation: &nav.Location{Lat: 2.2, Lon: 2.2},
					Distance:        11.1,
					Duration:        11.1,
				},
				{
					StationID:       testSGStationID3,
					StationLocation: &nav.Location{Lat: 3.3, Lon: 3.3},
					Distance:        33.3,
					Duration:        33.3,
				},
			},
			testSGStationID1: {
				{
					StationID:       testSGStationID4,
					StationLocation: &nav.Location{Lat: 4.4, Lon: 4.4},
					Distance:        44.4,
					Duration:        44.4,
				},
				{
					StationID:       testSGStationID5,
					StationLocation: &nav.Location{Lat: 5.5, Lon: 5.5},
					Distance:        34.4,
					Duration:        34.4,
				},
			},
			testSGStationID2: {
				{
					StationID:       testSGStationID4,
					StationLocation: &nav.Location{Lat: 4.4, Lon: 4.4},
					Distance:        11.1,
					Duration:        11.1,
				},
				{
					StationID:       testSGStationID5,
					StationLocation: &nav.Location{Lat: 5.5, Lon: 5.5},
					Distance:        14.4,
					Duration:        14.4,
				},
			},
			testSGStationID3: {
				{
					StationID:       testSGStationID4,
					StationLocation: &nav.Location{Lat: 4.4, Lon: 4.4},
					Distance:        22.2,
					Duration:        22.2,
				},
				{
					StationID:       testSGStationID5,
					StationLocation: &nav.Location{Lat: 5.5, Lon: 5.5},
					Distance:        15.5,
					Duration:        15.5,
				},
			},
			testSGStationID4: {
				{
					StationID:       stationfindertype.DestLocationID,
					StationLocation: &nav.Location{Lat: 6.6, Lon: 6.6},
					Distance:        44.4,
					Duration:        44.4,
				},
			},
			testSGStationID5: {
				{
					StationID:       stationfindertype.DestLocationID,
					StationLocation: &nav.Location{Lat: 6.6, Lon: 6.6},
					Distance:        33.3,
					Duration:        33.3,
				},
			},
			stationfindertype.DestLocationID: {},
		},
		mockStationID2Location: map[string]*nav.Location{
			stationfindertype.OrigLocationID: {Lat: 0.0, Lon: 0.0},
			testSGStationID1:                 {Lat: 1.1, Lon: 1.1},
			testSGStationID2:                 {Lat: 2.2, Lon: 2.2},
			testSGStationID3:                 {Lat: 3.3, Lon: 3.3},
			testSGStationID4:                 {Lat: 4.4, Lon: 4.4},
			testSGStationID5:                 {Lat: 5.5, Lon: 5.5},
			stationfindertype.DestLocationID: {Lat: 6.6, Lon: 6.6},
		},
	}

	return querier
}

func (querier *mockQuerier4StationGraph) NearByStationQuery(stationID string) []*connectivitymap.QueryResult {
	if queryResult, ok := querier.mockStationID2QueryResult[stationID]; ok {
		return queryResult
	}
	glog.Fatal("Un-implemented mapping key for mockStationID2QueryResult.\n")
	return nil
}

func (querier *mockQuerier4StationGraph) GetLocation(stationID string) *nav.Location {
	if location, ok := querier.mockStationID2Location[stationID]; ok {
		return location
	}
	glog.Fatal("Un-implemented mapping key for mockStationID2Location.\n")
	return nil
}

func TestStationGraphGenerateSolutions1(t *testing.T) {
	maxEnergyLevel := 50.0
	currEnergyLevel := 20.0
	strategy := chargingstrategy.NewFakeChargingStrategy(maxEnergyLevel)
	querier := newMockQuerier4StationGraph()

	solutions := NewStationGraph(currEnergyLevel, maxEnergyLevel, strategy, querier).GenerateChargeSolutions()

	if len(solutions) != 0 {
		t.Errorf("expect to have 1 solution but got %d.\n", len(solutions))
	}
	// 	sol := solutions[0]
	// 	// 58.8 = 11.1 + 14.4 + 33.3
	// 	if !util.FloatEquals(sol.Distance, 58.8) {
	// 		t.Errorf("Incorrect distance calculated for fakeGraph1 expect 58.89 but got %#v.\n", sol.Distance)
	// 	}

	// 	// 7918.8 = 11.1 + 2532(60% charge) + 14.4 + 5328(80% charge) + 33.3
	// 	if !util.FloatEquals(sol.Duration, 7918.8) {
	// 		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.Duration)
	// 	}

	// 	// 6.7 = 40 - 33.3
	// 	if !util.FloatEquals(sol.RemainingRage, 6.7) {
	// 		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.RemainingRage)
	// 	}

	// 	if len(sol.ChargeStations) != 2 {
	// 		t.Errorf("Expect to have 2 charge stations for fakeGraph1 but got %d.\n", len(sol.ChargeStations))
	// 	}

	// 	expectStation1 := &solution.ChargeStation{
	// 		Location: nav.Location{
	// 			Lat: 2.2,
	// 			Lon: 2.2,
	// 		},
	// 		StationID:     "station2",
	// 		ArrivalEnergy: 8.9,
	// 		WaitTime:      0,
	// 		ChargeTime:    2532,
	// 		ChargeRange:   30,
	// 	}
	// 	if !reflect.DeepEqual(sol.ChargeStations[0], expectStation1) {
	// 		t.Errorf("Expect first charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation1, sol.ChargeStations[0])
	// 	}

	// 	expectStation2 := &solution.ChargeStation{
	// 		Location: nav.Location{
	// 			Lat: 5.5,
	// 			Lon: 5.5,
	// 		},
	// 		StationID:     "station5",
	// 		ArrivalEnergy: 15.6,
	// 		WaitTime:      0,
	// 		ChargeTime:    5328,
	// 		ChargeRange:   40,
	// 	}
	// 	if !reflect.DeepEqual(sol.ChargeStations[1], expectStation2) {
	// 		t.Errorf("Expect second charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation2, sol.ChargeStations[1])
	// 	}
}

var mockedGraph = mockGraph{
	[]*node{
		{
			0,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 20.0,
				},
			},
			locationInfo{
				lat: 0.0,
				lon: 0.0,
			},
		},
		{
			1,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			locationInfo{
				lat: 1.1,
				lon: 1.1,
			},
		},
		{
			2,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			locationInfo{
				lat: 1.1,
				lon: 1.1,
			},
		},
		{
			3,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			locationInfo{
				lat: 1.1,
				lon: 1.1,
			},
		},
		{
			4,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			5,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			6,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			locationInfo{
				lat: 2.2,
				lon: 2.2,
			},
		},
		{
			7,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			8,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			9,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			locationInfo{
				lat: 3.3,
				lon: 3.3,
			},
		},
		{
			10,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
		},
		{
			11,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
		},
		{
			12,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			locationInfo{
				lat: 4.4,
				lon: 4.4,
			},
		},
		{
			13,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			locationInfo{
				lat: 5.5,
				lon: 5.5,
			},
		},
		{
			14,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			locationInfo{
				lat: 5.5,
				lon: 5.5,
			},
		},
		{
			15,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			locationInfo{
				lat: 5.5,
				lon: 5.5,
			},
		},
		{
			16,
			chargeInfo{
				targetState: chargingstrategy.State{},
			},
			locationInfo{
				lat: 6.6,
				lon: 6.6,
			},
		},
	},
	[]string{
		stationfindertype.OrigLocationID,
		"station1",
		"station1",
		"station1",
		"station2",
		"station2",
		"station2",
		"station3",
		"station3",
		"station3",
		"station4",
		"station4",
		"station4",
		"station5",
		"station5",
		"station5",
		stationfindertype.DestLocationID,
	},
	map[nodeID][]*edgeIDAndData{
		0: {
			// orig -> station 1
			{edgeID{0, 1}, &edge{22.2, 22.2}},
			{edgeID{0, 2}, &edge{22.2, 22.2}},
			{edgeID{0, 3}, &edge{22.2, 22.2}},
			// orig -> station 2
			{edgeID{0, 4}, &edge{11.1, 11.1}},
			{edgeID{0, 5}, &edge{11.1, 11.1}},
			{edgeID{0, 6}, &edge{11.1, 11.1}},
			// orig -> station 3
			{edgeID{0, 7}, &edge{33.3, 33.3}},
			{edgeID{0, 8}, &edge{33.3, 33.3}},
			{edgeID{0, 9}, &edge{33.3, 33.3}},
		},
		1: {
			// station 1 -> station 4
			{edgeID{1, 10}, &edge{44.4, 44.4}},
			{edgeID{1, 11}, &edge{44.4, 44.4}},
			{edgeID{1, 12}, &edge{44.4, 44.4}},
			// station 1 -> station 5
			{edgeID{1, 13}, &edge{34.4, 34.4}},
			{edgeID{1, 14}, &edge{34.4, 34.4}},
			{edgeID{1, 15}, &edge{34.4, 34.4}},
		},
		2: {
			// station 1 -> station 4
			{edgeID{2, 10}, &edge{44.4, 44.4}},
			{edgeID{2, 11}, &edge{44.4, 44.4}},
			{edgeID{2, 12}, &edge{44.4, 44.4}},
			// station 1 -> station 5
			{edgeID{2, 13}, &edge{34.4, 34.4}},
			{edgeID{2, 14}, &edge{34.4, 34.4}},
			{edgeID{2, 15}, &edge{34.4, 34.4}},
		},
		3: {
			// station 1 -> station 4
			{edgeID{3, 10}, &edge{44.4, 44.4}},
			{edgeID{3, 11}, &edge{44.4, 44.4}},
			{edgeID{3, 12}, &edge{44.4, 44.4}},
			// station 1 -> station 5
			{edgeID{3, 13}, &edge{34.4, 34.4}},
			{edgeID{3, 14}, &edge{34.4, 34.4}},
			{edgeID{3, 15}, &edge{34.4, 34.4}},
		},
		4: {
			// station 2 -> station 4
			{edgeID{4, 10}, &edge{11.1, 11.1}},
			{edgeID{4, 11}, &edge{11.1, 11.1}},
			{edgeID{4, 12}, &edge{11.1, 11.1}},
			// station 2 -> station 5
			{edgeID{4, 13}, &edge{14.4, 14.4}},
			{edgeID{4, 14}, &edge{14.4, 14.4}},
			{edgeID{4, 15}, &edge{14.4, 14.4}},
		},
		5: {
			// station 2 -> station 4
			{edgeID{5, 10}, &edge{11.1, 11.1}},
			{edgeID{5, 11}, &edge{11.1, 11.1}},
			{edgeID{5, 12}, &edge{11.1, 11.1}},
			// station 2 -> station 5
			{edgeID{5, 13}, &edge{14.4, 14.4}},
			{edgeID{5, 14}, &edge{14.4, 14.4}},
			{edgeID{5, 15}, &edge{14.4, 14.4}},
		},
		6: {
			// station 2 -> station 4
			{edgeID{6, 10}, &edge{11.1, 11.1}},
			{edgeID{6, 11}, &edge{11.1, 11.1}},
			{edgeID{6, 12}, &edge{11.1, 11.1}},
			// station 2 -> station 5
			{edgeID{6, 13}, &edge{14.4, 14.4}},
			{edgeID{6, 14}, &edge{14.4, 14.4}},
			{edgeID{6, 15}, &edge{14.4, 14.4}},
		},
		7: {
			// station 3 -> station 4
			{edgeID{7, 10}, &edge{22.2, 22.2}},
			{edgeID{7, 11}, &edge{22.2, 22.2}},
			{edgeID{7, 12}, &edge{22.2, 22.2}},
			// station 3 -> station 5
			{edgeID{7, 13}, &edge{15.5, 15.5}},
			{edgeID{7, 14}, &edge{15.5, 15.5}},
			{edgeID{7, 15}, &edge{15.5, 15.5}},
		},
		8: {
			// station 3 -> station 4
			{edgeID{8, 10}, &edge{22.2, 22.2}},
			{edgeID{8, 11}, &edge{22.2, 22.2}},
			{edgeID{8, 12}, &edge{22.2, 22.2}},
			// station 3 -> station 5
			{edgeID{8, 13}, &edge{15.5, 15.5}},
			{edgeID{8, 14}, &edge{15.5, 15.5}},
			{edgeID{8, 15}, &edge{15.5, 15.5}},
		},
		9: {
			// station 3 -> station 4
			{edgeID{9, 10}, &edge{22.2, 22.2}},
			{edgeID{9, 11}, &edge{22.2, 22.2}},
			{edgeID{9, 12}, &edge{22.2, 22.2}},
			// station 3 -> station 5
			{edgeID{9, 13}, &edge{15.5, 15.5}},
			{edgeID{9, 14}, &edge{15.5, 15.5}},
			{edgeID{9, 15}, &edge{15.5, 15.5}},
		},
		10: {
			// station 4 -> end
			{edgeID{10, 16}, &edge{44.4, 44.4}},
			{edgeID{10, 16}, &edge{44.4, 44.4}},
			{edgeID{10, 16}, &edge{44.4, 44.4}},
		},
		11: {
			// station 4 -> end
			{edgeID{11, 16}, &edge{44.4, 44.4}},
			{edgeID{11, 16}, &edge{44.4, 44.4}},
			{edgeID{11, 16}, &edge{44.4, 44.4}},
		},
		12: {
			// station 4 -> end
			{edgeID{12, 16}, &edge{44.4, 44.4}},
			{edgeID{12, 16}, &edge{44.4, 44.4}},
			{edgeID{12, 16}, &edge{44.4, 44.4}},
		},
		13: {
			// station 5 -> end
			{edgeID{13, 16}, &edge{33.3, 33.3}},
			{edgeID{13, 16}, &edge{33.3, 33.3}},
			{edgeID{13, 16}, &edge{33.3, 33.3}},
		},
		14: {
			// station 5 -> end
			{edgeID{14, 16}, &edge{33.3, 33.3}},
			{edgeID{14, 16}, &edge{33.3, 33.3}},
			{edgeID{14, 16}, &edge{33.3, 33.3}},
		},
		15: {
			// station 5 -> end
			{edgeID{15, 16}, &edge{33.3, 33.3}},
			{edgeID{15, 16}, &edge{33.3, 33.3}},
			{edgeID{15, 16}, &edge{33.3, 33.3}},
		},
	},
	chargingstrategy.NewFakeChargingStrategy(50.0),
}

// var fakeNeighborsGraph = [][]stationfindertype.NeighborInfo{
// 	[]stationfindertype.NeighborInfo{
// 		stationfindertype.NeighborInfo{
// 			FromID: "orig_location",
// 			FromLocation: nav.Location{
// 				Lat: 0.0,
// 				Lon: 0.0,
// 			},
// 			ToID: "station1",
// 			ToLocation: nav.Location{
// 				Lat: 1.1,
// 				Lon: 1.1,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 22.2,
// 				Distance: 22.2,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "orig_location",
// 			FromLocation: nav.Location{
// 				Lat: 0.0,
// 				Lon: 0.0,
// 			},
// 			ToID: "station2",
// 			ToLocation: nav.Location{
// 				Lat: 2.2,
// 				Lon: 2.2,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 11.1,
// 				Distance: 11.1,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "orig_location",
// 			FromLocation: nav.Location{
// 				Lat: 0.0,
// 				Lon: 0.0,
// 			},
// 			ToID: "station3",
// 			Weight: stationfindertype.Weight{
// 				Duration: 33.3,
// 				Distance: 33.3,
// 			},
// 			ToLocation: nav.Location{
// 				Lat: 3.3,
// 				Lon: 3.3,
// 			},
// 		},
// 	},
// 	[]stationfindertype.NeighborInfo{
// 		stationfindertype.NeighborInfo{
// 			FromID: "station1",
// 			FromLocation: nav.Location{
// 				Lat: 1.1,
// 				Lon: 1.1,
// 			},
// 			ToID: "station4",
// 			ToLocation: nav.Location{
// 				Lat: 4.4,
// 				Lon: 4.4,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 44.4,
// 				Distance: 44.4,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station1",
// 			FromLocation: nav.Location{
// 				Lat: 1.1,
// 				Lon: 1.1,
// 			},
// 			ToID: "station5",
// 			ToLocation: nav.Location{
// 				Lat: 5.5,
// 				Lon: 5.5,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 34.4,
// 				Distance: 34.4,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station2",
// 			FromLocation: nav.Location{
// 				Lat: 2.2,
// 				Lon: 2.2,
// 			},
// 			ToID: "station4",
// 			ToLocation: nav.Location{
// 				Lat: 4.4,
// 				Lon: 4.4,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 11.1,
// 				Distance: 11.1,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station2",
// 			FromLocation: nav.Location{
// 				Lat: 2.2,
// 				Lon: 2.2,
// 			},
// 			ToID: "station5",
// 			ToLocation: nav.Location{
// 				Lat: 5.5,
// 				Lon: 5.5,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 14.4,
// 				Distance: 14.4,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station3",
// 			FromLocation: nav.Location{
// 				Lat: 3.3,
// 				Lon: 3.3,
// 			},
// 			ToID: "station4",
// 			ToLocation: nav.Location{
// 				Lat: 4.4,
// 				Lon: 4.4,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 22.2,
// 				Distance: 22.2,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station3",
// 			FromLocation: nav.Location{
// 				Lat: 3.3,
// 				Lon: 3.3,
// 			},
// 			ToID: "station5",
// 			ToLocation: nav.Location{
// 				Lat: 5.5,
// 				Lon: 5.5,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 15.5,
// 				Distance: 15.5,
// 			},
// 		},
// 	},
// 	[]stationfindertype.NeighborInfo{
// 		stationfindertype.NeighborInfo{
// 			FromID: "station4",
// 			FromLocation: nav.Location{
// 				Lat: 4.4,
// 				Lon: 4.4,
// 			},
// 			ToID: stationfindertype.DestLocationID,
// 			ToLocation: nav.Location{
// 				Lat: 6.6,
// 				Lon: 6.6,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 44.4,
// 				Distance: 44.4,
// 			},
// 		},
// 		stationfindertype.NeighborInfo{
// 			FromID: "station5",
// 			FromLocation: nav.Location{
// 				Lat: 5.5,
// 				Lon: 5.5,
// 			},
// 			ToID: stationfindertype.DestLocationID,
// 			ToLocation: nav.Location{
// 				Lat: 6.6,
// 				Lon: 6.6,
// 			},
// 			Weight: stationfindertype.Weight{
// 				Duration: 33.3,
// 				Distance: 33.3,
// 			},
// 		},
// 	},
// }

// func TestConstructStationGraph(t *testing.T) {
// 	// generate channel contains neighbors information
// 	// simulate real situation using different go-routine
// 	c := make(chan stationfindertype.WeightBetweenNeighbors)
// 	go func() {
// 		for _, n := range fakeNeighborsGraph {
// 			neighborsInfo := stationfindertype.WeightBetweenNeighbors{
// 				NeighborsInfo: n,
// 				Err:           nil,
// 			}
// 			c <- neighborsInfo
// 		}
// 		close(c)
// 	}()

// 	currEnergyLevel := 0.0
// 	maxEnergyLevel := 50.0
// 	graph := NewStationGraph(c, currEnergyLevel, maxEnergyLevel,
// 		chargingstrategy.NewFakeChargingStrategy(maxEnergyLevel))
// 	if graph == nil {
// 		t.Errorf("create Station graph failed, expect none-empty graph but result is empty")
// 	}

// 	testStart(t, graph, currEnergyLevel, maxEnergyLevel)
// 	testEnd(t, graph, currEnergyLevel, maxEnergyLevel)

// 	testConnectivity(t, graph, "station1", locationInfo{lat: 1.1, lon: 1.1},
// 		[]string{"station4", "station5"}, fakeNeighborsGraph, currEnergyLevel, maxEnergyLevel)

// 	testConnectivity(t, graph, "station2", locationInfo{lat: 2.2, lon: 2.2},
// 		[]string{"station4", "station5"}, fakeNeighborsGraph, currEnergyLevel, maxEnergyLevel)

// 	testConnectivity(t, graph, "station3", locationInfo{lat: 3.3, lon: 3.3},
// 		[]string{"station4", "station5"}, fakeNeighborsGraph, currEnergyLevel, maxEnergyLevel)

// 	testConnectivity(t, graph, "station4", locationInfo{lat: 4.4, lon: 4.4},
// 		[]string{stationfindertype.DestLocationID}, fakeNeighborsGraph, currEnergyLevel, maxEnergyLevel)

// 	testConnectivity(t, graph, "station5", locationInfo{lat: 5.5, lon: 5.5},
// 		[]string{stationfindertype.DestLocationID}, fakeNeighborsGraph, currEnergyLevel, maxEnergyLevel)
// }

// func testStart(t *testing.T, graph *stationGraph, currEnergyLevel, maxEnergyLevel float64) {
// 	sn := graph.getChargeStationsNodes(stationfindertype.OrigLocationID, nav.Location{}, currEnergyLevel, maxEnergyLevel)
// 	if len(sn) != 1 {
// 		t.Errorf("incorrect start node generated expect only one node but got %d", len(sn))
// 	}
// 	if graph.getStationID(sn[0].id) != stationfindertype.OrigLocationID {
// 		t.Errorf("incorrect name for start node expect %s but got %s", stationfindertype.OrigLocationID, graph.getStationID(sn[0].id))
// 	}
// 	if !util.FloatEquals(sn[0].arrivalEnergy, currEnergyLevel) ||
// 		!util.FloatEquals(sn[0].targetState.Energy, 0.0) ||
// 		!util.FloatEquals(sn[0].chargeTime, 0.0) {
// 		t.Errorf("incorrect energy information for start node expect %v but got %v", chargeInfo{
// 			arrivalEnergy: currEnergyLevel,
// 			chargeTime:    0.0,
// 			targetState: chargingstrategy.State{
// 				Energy: 0.0,
// 			},
// 		}, sn[0].chargeInfo)
// 	}

// 	startLocation := graph.g.getLocationInfo(sn[0].id)
// 	if !util.FloatEquals(startLocation.lat, 0.0) ||
// 		!util.FloatEquals(startLocation.lon, 0.0) {
// 		t.Errorf("incorrect location information for start node expect %v but got %v", locationInfo{
// 			lat: 0.0,
// 			lon: 0.0,
// 		}, startLocation)
// 	}

// 	if len(sn[0].neighbors) != 9 {
// 		t.Errorf("incorrect neighbors count for start node expect %d but got %d", 9, len(sn[0].neighbors))
// 	}

// 	if graph.getStationID(sn[0].neighbors[0].targetNodeID) != "station1" ||
// 		!util.FloatEquals(sn[0].neighbors[0].distance, 22.2) ||
// 		!util.FloatEquals(sn[0].neighbors[0].duration, 22.2) ||
// 		graph.getStationID(sn[0].neighbors[1].targetNodeID) != "station1" ||
// 		!util.FloatEquals(sn[0].neighbors[1].distance, 22.2) ||
// 		!util.FloatEquals(sn[0].neighbors[1].duration, 22.2) ||
// 		graph.getStationID(sn[0].neighbors[2].targetNodeID) != "station1" ||
// 		!util.FloatEquals(sn[0].neighbors[2].distance, 22.2) ||
// 		!util.FloatEquals(sn[0].neighbors[2].duration, 22.2) ||
// 		graph.getStationID(sn[0].neighbors[3].targetNodeID) != "station2" ||
// 		!util.FloatEquals(sn[0].neighbors[3].distance, 11.1) ||
// 		!util.FloatEquals(sn[0].neighbors[3].duration, 11.1) ||
// 		graph.getStationID(sn[0].neighbors[4].targetNodeID) != "station2" ||
// 		!util.FloatEquals(sn[0].neighbors[4].distance, 11.1) ||
// 		!util.FloatEquals(sn[0].neighbors[4].duration, 11.1) ||
// 		graph.getStationID(sn[0].neighbors[5].targetNodeID) != "station2" ||
// 		!util.FloatEquals(sn[0].neighbors[5].distance, 11.1) ||
// 		!util.FloatEquals(sn[0].neighbors[5].duration, 11.1) ||
// 		graph.getStationID(sn[0].neighbors[6].targetNodeID) != "station3" ||
// 		!util.FloatEquals(sn[0].neighbors[6].distance, 33.3) ||
// 		!util.FloatEquals(sn[0].neighbors[6].duration, 33.3) ||
// 		graph.getStationID(sn[0].neighbors[7].targetNodeID) != "station3" ||
// 		!util.FloatEquals(sn[0].neighbors[7].distance, 33.3) ||
// 		!util.FloatEquals(sn[0].neighbors[7].duration, 33.3) ||
// 		graph.getStationID(sn[0].neighbors[8].targetNodeID) != "station3" ||
// 		!util.FloatEquals(sn[0].neighbors[8].distance, 33.3) ||
// 		!util.FloatEquals(sn[0].neighbors[8].duration, 33.3) {
// 		t.Errorf("incorrect neighbor information generated for start node")
// 	}
// }

// func testEnd(t *testing.T, graph *stationGraph, currEnergyLevel, maxEnergyLevel float64) {
// 	se := graph.getChargeStationsNodes(stationfindertype.DestLocationID, nav.Location{}, currEnergyLevel, maxEnergyLevel)
// 	if len(se) != 1 {
// 		t.Errorf("incorrect end node generated expect only one node but got %d", len(se))
// 	}
// 	if graph.getStationID(se[0].id) != stationfindertype.DestLocationID {
// 		t.Errorf("incorrect name for end node expect %s but got %s", stationfindertype.DestLocationID, graph.getStationID(se[0].id))
// 	}
// 	if !util.FloatEquals(se[0].arrivalEnergy, 0.0) ||
// 		!util.FloatEquals(se[0].targetState.Energy, 0.0) ||
// 		!util.FloatEquals(se[0].chargeTime, 0.0) {
// 		t.Errorf("incorrect energy information for end node expect %v but got %v", chargeInfo{
// 			arrivalEnergy: 0.0,
// 			chargeTime:    0.0,
// 			targetState: chargingstrategy.State{
// 				Energy: 0.0,
// 			},
// 		}, se[0].chargeInfo)
// 	}

// 	endLocation := graph.g.getLocationInfo(se[0].id)
// 	if !util.FloatEquals(endLocation.lat, 6.6) ||
// 		!util.FloatEquals(endLocation.lon, 6.6) {
// 		t.Errorf("incorrect location information for end node expect %v but got %v", locationInfo{
// 			lat: 6.6,
// 			lon: 6.6,
// 		}, endLocation)
// 	}

// 	// if len(se[0].neighbors) != 0 {
// 	// 	t.Errorf("incorrect neighbors count for end node expect %d but got %d", 0, len(se[0].neighbors))
// 	// }
// }

// func testConnectivity(t *testing.T, graph *stationGraph, from string, fromLocation locationInfo,
// 	tos []string, mockArray [][]stationfindertype.NeighborInfo, currEnergyLevel, maxEnergyLevel float64) {
// 	fns := graph.getChargeStationsNodes(from, nav.Location{}, 0.0, 0.0)

// 	for _, fromNode := range fns {
// 		if !util.FloatEquals(fromNode.locationInfo.lat, fromLocation.lat) ||
// 			!util.FloatEquals(fromNode.locationInfo.lon, fromLocation.lon) {
// 			t.Errorf("incorrect location information generated for node %s expect %+v got %+v",
// 				from, fromLocation, fromNode.locationInfo)
// 		}
// 	}

// 	index := 0
// 	for _, to := range tos {
// 		tns := graph.getChargeStationsNodes(to, nav.Location{}, 0.0, 0.0)

// 		expectDuration := math.MaxFloat64
// 		expectDistance := math.MaxFloat64
// 		for _, neighborsInfo := range mockArray {
// 			for _, neighborInfo := range neighborsInfo {
// 				if neighborInfo.FromID == from && neighborInfo.ToID == to {
// 					expectDuration = neighborInfo.Duration
// 					expectDistance = neighborInfo.Distance
// 					break
// 				}
// 			}
// 		}
// 		if expectDuration == math.MaxFloat64 ||
// 			expectDistance == math.MaxFloat64 {
// 			t.Error("incorrect name string passed into testConnectivity")
// 		}

// 		for _, fromNode := range fns {
// 			for i, toNode := range tns {
// 				if fromNode.neighbors[index+i].targetNodeID != toNode.id ||
// 					fromNode.neighbors[index+i].distance != expectDistance ||
// 					fromNode.neighbors[index+i].duration != expectDuration {
// 					t.Errorf("incorrect connectivity generated between %s and %s", from, to)
// 				}
// 			}
// 		}

// 		index += len(tns)
// 	}
// }

// // based on original graph, best charge solution is
// // start -> station 2 -> station 4 -> end
// // when start, initial energy is 20
// // start -> station 2, time/duration = 11.1, this case will choose charging for 60%
// // station 2 -> station 5, time/duration = 14.4, this cause will choose charging for 80%
// func TestGenerateChargeSolutions1(t *testing.T) {

// 	fakeGraph1 := make([][]stationfindertype.NeighborInfo, len(fakeNeighborsGraph))
// 	for i := range fakeNeighborsGraph {
// 		fakeGraph1[i] = make([]stationfindertype.NeighborInfo, len(fakeNeighborsGraph[i]))
// 		copy(fakeGraph1[i], fakeNeighborsGraph[i])
// 	}

// 	// generate channel contains neighbors information
// 	// simulate real situation using different go-routine
// 	c := make(chan stationfindertype.WeightBetweenNeighbors)
// 	go func() {
// 		for _, n := range fakeGraph1 {
// 			neighborsInfo := stationfindertype.WeightBetweenNeighbors{
// 				NeighborsInfo: n,
// 				Err:           nil,
// 			}
// 			c <- neighborsInfo
// 		}
// 		close(c)
// 	}()

// 	currEnergyLevel := 20.0
// 	maxEnergyLevel := 50.0
// 	graph := NewStationGraph(c, currEnergyLevel, maxEnergyLevel,
// 		chargingstrategy.NewFakeChargingStrategy(maxEnergyLevel))
// 	if graph == nil {
// 		t.Error("create Station graph failed, expect none-empty graph but result is empty")
// 	}

// 	solutions := graph.GenerateChargeSolutions()
// 	fmt.Printf("### %#v\n", solutions[0])
// 	fmt.Printf("### %#v\n", solutions[0].ChargeStations[0])
// 	fmt.Printf("### %#v\n", solutions[0].ChargeStations[1])
// 	if len(solutions) != 1 {
// 		t.Errorf("expect to have 1 solution but got %d.\n", len(solutions))
// 	}
// 	sol := solutions[0]
// 	// 58.8 = 11.1 + 14.4 + 33.3
// 	if !util.FloatEquals(sol.Distance, 58.8) {
// 		t.Errorf("Incorrect distance calculated for fakeGraph1 expect 58.89 but got %#v.\n", sol.Distance)
// 	}

// 	// 7918.8 = 11.1 + 2532(60% charge) + 14.4 + 5328(80% charge) + 33.3
// 	if !util.FloatEquals(sol.Duration, 7918.8) {
// 		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.Duration)
// 	}

// 	// 6.7 = 40 - 33.3
// 	if !util.FloatEquals(sol.RemainingRage, 6.7) {
// 		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.RemainingRage)
// 	}

// 	if len(sol.ChargeStations) != 2 {
// 		t.Errorf("Expect to have 2 charge stations for fakeGraph1 but got %d.\n", len(sol.ChargeStations))
// 	}

// 	expectStation1 := &solution.ChargeStation{
// 		Location: nav.Location{
// 			Lat: 2.2,
// 			Lon: 2.2,
// 		},
// 		StationID:     "station2",
// 		ArrivalEnergy: 8.9,
// 		WaitTime:      0,
// 		ChargeTime:    2532,
// 		ChargeRange:   30,
// 	}
// 	if !reflect.DeepEqual(sol.ChargeStations[0], expectStation1) {
// 		t.Errorf("Expect first charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation1, sol.ChargeStations[0])
// 	}

// 	expectStation2 := &solution.ChargeStation{
// 		Location: nav.Location{
// 			Lat: 5.5,
// 			Lon: 5.5,
// 		},
// 		StationID:     "station5",
// 		ArrivalEnergy: 15.6,
// 		WaitTime:      0,
// 		ChargeTime:    5328,
// 		ChargeRange:   40,
// 	}
// 	if !reflect.DeepEqual(sol.ChargeStations[1], expectStation2) {
// 		t.Errorf("Expect second charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation2, sol.ChargeStations[1])
// 	}

// }
