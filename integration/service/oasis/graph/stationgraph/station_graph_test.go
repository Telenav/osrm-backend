package stationgraph

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/Telenav/osrm-backend/integration/util"
	"github.com/golang/glog"
)

/*
Construct test graph of
- start connects to staion 1, station 2, station 3
     - start -> station 1: 22.2,
     - start -> station 2: 11.1,
     - start -> station 3: 33.3,
- station 1 connects to station 4, station 5,
     - station 1 -> station 4: 44.4,
     - station 1 -> station 5: 34.4,
- station 2 connects to station 4, station 5
     - station 2 -> station 4: 11.1,
     - station 2 -> station 5: 14.4,
- station 3 connects to station 4, station 5
     - station 3 -> station 4: 22.2,
     - station 3 -> station 5: 15.5,
- station 4 connects to end
     - station 4 -> end      : 44.4,
- station 5 connects to end
     - station 5 -> end      : 33.3,
                    station 1
               /       \    \
              /         \    \
             /          _\_____station 4
            /          /  \     /        \
           /          /    \   /          \
start  -------   station 2  \ /           end
           \          \     / \          /
            \          \   /   \        /
             \          \_/_____station 5
              \          /     /
               \        /     /
                    station 3

- Each charge station will try to generate up to 3 virtual node based on fakechargestrategy,
  each node represent for situation of time spend in charging and get to different energy level
- In total there could be 17 nodes in graph
	  + start node
	  + end node
	  + station 1 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
	  + station 2 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
	  + station 3 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
	  + station 4 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
	  + station 5 * 3 (node represent for with 60% of total energy, with 80% of total energy, with 100% of total energy, respectively)
- Take station 1 as an example, it will create 3 nodes. For each of them, will contains neighbor information with:
	  + connection to station 4's node with 60% of total energy
	  + connection to station 4's node with 80% of total energy
	  + connection to station 4's node with 100% of total energy
	  + connection to station 5's node with 60% of total energy
	  + connection to station 5's node with 80% of total energy
	  + connection to station 5's node with 100% of total energy
  Each node with name `station1` will have different id, and each of them will have 6 neighbor nodes
- For default graph, initial energy is 20.0, max energy is 50.0
*/

var testSGStationID1Str = "1"
var testSGStationID2Str = "2"
var testSGStationID3Str = "3"
var testSGStationID4Str = "4"
var testSGStationID5Str = "5"

var testSGStationID1 entity.PlaceID = 1
var testSGStationID2 entity.PlaceID = 2
var testSGStationID3 entity.PlaceID = 3
var testSGStationID4 entity.PlaceID = 4
var testSGStationID5 entity.PlaceID = 5

type mockQuerier4StationGraph struct {
	mockStationID2QueryResult map[entity.PlaceID][]*entity.TransferInfo
	mockStationID2Location    map[entity.PlaceID]*nav.Location
}

func newMockQuerier4StationGraph() place.TopoQuerier {
	querier := &mockQuerier4StationGraph{
		mockStationID2QueryResult: map[entity.PlaceID][]*entity.TransferInfo{
			iteratortype.OrigLocationID: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID2,
						Location: &nav.Location{Lat: 2.2, Lon: 2.2},
					},
					Weight: &entity.Weight{
						Distance: 11.1,
						Duration: 11.1,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID1,
						Location: &nav.Location{Lat: 1.1, Lon: 1.1},
					},
					Weight: &entity.Weight{
						Distance: 22.2,
						Duration: 22.2,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID3,
						Location: &nav.Location{Lat: 3.3, Lon: 3.3},
					},
					Weight: &entity.Weight{
						Distance: 33.3,
						Duration: 33.3,
					},
				},
			},
			testSGStationID1: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID5,
						Location: &nav.Location{Lat: 5.5, Lon: 5.5},
					},
					Weight: &entity.Weight{
						Distance: 34.4,
						Duration: 34.4,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID4,
						Location: &nav.Location{Lat: 4.4, Lon: 4.4},
					},
					Weight: &entity.Weight{
						Distance: 44.4,
						Duration: 44.4,
					},
				},
			},
			testSGStationID2: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID4,
						Location: &nav.Location{Lat: 4.4, Lon: 4.4},
					},
					Weight: &entity.Weight{
						Distance: 11.1,
						Duration: 11.1,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID5,
						Location: &nav.Location{Lat: 5.5, Lon: 5.5},
					},
					Weight: &entity.Weight{
						Distance: 14.4,
						Duration: 14.4,
					},
				},
			},
			testSGStationID3: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID5,
						Location: &nav.Location{Lat: 5.5, Lon: 5.5},
					},
					Weight: &entity.Weight{
						Distance: 15.5,
						Duration: 15.5,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       testSGStationID4,
						Location: &nav.Location{Lat: 4.4, Lon: 4.4},
					},
					Weight: &entity.Weight{
						Distance: 22.2,
						Duration: 22.2,
					},
				},
			},
			testSGStationID4: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       iteratortype.DestLocationID,
						Location: &nav.Location{Lat: 6.6, Lon: 6.6},
					},
					Weight: &entity.Weight{
						Distance: 44.4,
						Duration: 44.4,
					},
				},
			},
			testSGStationID5: {
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID:       iteratortype.DestLocationID,
						Location: &nav.Location{Lat: 6.6, Lon: 6.6},
					},
					Weight: &entity.Weight{
						Distance: 33.3,
						Duration: 33.3,
					},
				},
			},
			iteratortype.DestLocationID: {},
		},
		mockStationID2Location: map[entity.PlaceID]*nav.Location{
			iteratortype.OrigLocationID: {Lat: 0.0, Lon: 0.0},
			testSGStationID1:            {Lat: 1.1, Lon: 1.1},
			testSGStationID2:            {Lat: 2.2, Lon: 2.2},
			testSGStationID3:            {Lat: 3.3, Lon: 3.3},
			testSGStationID4:            {Lat: 4.4, Lon: 4.4},
			testSGStationID5:            {Lat: 5.5, Lon: 5.5},
			iteratortype.DestLocationID: {Lat: 6.6, Lon: 6.6},
		},
	}

	return querier
}

func (querier *mockQuerier4StationGraph) GetConnectedPlaces(placeID entity.PlaceID) []*entity.TransferInfo {
	if queryResult, ok := querier.mockStationID2QueryResult[placeID]; ok {
		return queryResult
	}
	glog.Fatalf("Un-implemented mapping key %v for mockStationID2QueryResult.\n", placeID)
	return nil
}

func (querier *mockQuerier4StationGraph) GetLocation(placeID entity.PlaceID) *nav.Location {
	if location, ok := querier.mockStationID2Location[placeID]; ok {
		return location
	}
	glog.Fatalf("Un-implemented mapping key for mockStationID2Location id = %v.\n", placeID)
	return nil
}

func TestStationGraphGenerateSolutions1(t *testing.T) {
	maxEnergyLevel := 50.0
	currEnergyLevel := 20.0
	strategy := chargingstrategy.NewSimpleChargingStrategy(maxEnergyLevel)
	querier := newMockQuerier4StationGraph()

	solutions := NewStationGraph(currEnergyLevel, maxEnergyLevel, strategy, querier).GenerateChargeSolutions()

	if len(solutions) != 1 {
		t.Errorf("expect to have 1 solution but got %d.\n", len(solutions))
	}
	sol := solutions[0]
	// 58.8 = 11.1 + 14.4 + 33.3
	if !util.Float64Equal(sol.Distance, 58.8) {
		t.Errorf("Incorrect distance calculated for fakeGraph1 expect 58.89 but got %#v.\n", sol.Distance)
	}

	// 7918.8 = 11.1 + 2532(60% charge) + 14.4 + 5328(80% charge) + 33.3
	if !util.Float64Equal(sol.Duration, 7918.8) {
		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.Duration)
	}

	// 6.7 = 40 - 33.3
	if !util.Float64Equal(sol.RemainingRage, 6.7) {
		t.Errorf("Incorrect duration calculated for fakeGraph1 expect 10858.8 but got %#v.\n", sol.RemainingRage)
	}

	if len(sol.ChargeStations) != 2 {
		t.Errorf("Expect to have 2 charge stations for fakeGraph1 but got %d.\n", len(sol.ChargeStations))
	}

	expectStation1 := &entity.ChargeStation{
		Location: nav.Location{
			Lat: 2.2,
			Lon: 2.2,
		},
		PlaceID:       "2",
		ArrivalEnergy: 8.9,
		WaitTime:      0,
		ChargeTime:    2532,
		ChargeRange:   30,
	}
	if !reflect.DeepEqual(sol.ChargeStations[0], expectStation1) {
		t.Errorf("Expect first charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation1, sol.ChargeStations[0])
	}

	expectStation2 := &entity.ChargeStation{
		Location: nav.Location{
			Lat: 5.5,
			Lon: 5.5,
		},
		PlaceID:       "5",
		ArrivalEnergy: 15.6,
		WaitTime:      0,
		ChargeTime:    5328,
		ChargeRange:   40,
	}
	if !reflect.DeepEqual(sol.ChargeStations[1], expectStation2) {
		t.Errorf("Expect second charge stations info for fakeGraph1 is %#v but got %#v\n", expectStation2, sol.ChargeStations[1])
	}
}

// mockedGraph4StationGraph defines compatible topological representation after mockQuerier4StationGraph is built
var mockedGraph4StationGraph = mockGraph{
	[]*node{
		// iteratortype.OrigLocationID.String(),
		{
			0,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 20.0,
				},
			},
			// nav.Location{
			// 	Lat: 0.0,
			// 	Lon: 0.0,
			// },
		},
		// iteratortype.DestLocationID.String(),
		{
			1,
			chargeInfo{
				targetState: chargingstrategy.State{},
			},
			// nav.Location{
			// 	Lat: 6.6,
			// 	Lon: 6.6,
			// },
		},
		// station 1
		{
			2,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			// nav.Location{
			// 	Lat: 1.1,
			// 	Lon: 1.1,
			// },
		},
		// station 1
		{
			3,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			// nav.Location{
			// 	Lat: 1.1,
			// 	Lon: 1.1,
			// },
		},
		// station 1
		{
			4,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			// nav.Location{
			// 	Lat: 1.1,
			// 	Lon: 1.1,
			// },
		},
		// station 2
		{
			5,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		// station 2
		{
			6,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		// station 2
		{
			7,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			// nav.Location{
			// 	Lat: 2.2,
			// 	Lon: 2.2,
			// },
		},
		// station 3
		{
			8,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		// station 3
		{
			9,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		// station 3
		{
			10,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			// nav.Location{
			// 	Lat: 3.3,
			// 	Lon: 3.3,
			// },
		},
		// station 4
		{
			11,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
		// station 4
		{
			12,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
		// station 4
		{
			13,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			// nav.Location{
			// 	Lat: 4.4,
			// 	Lon: 4.4,
			// },
		},
		// station 5
		{
			14,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 30.0,
				},
			},
			// nav.Location{
			// 	Lat: 5.5,
			// 	Lon: 5.5,
			// },
		},
		// station 5
		{
			15,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 40.0,
				},
			},
			// nav.Location{
			// 	Lat: 5.5,
			// 	Lon: 5.5,
			// },
		},
		// station 5
		{
			16,
			chargeInfo{
				targetState: chargingstrategy.State{
					Energy: 50.0,
				},
			},
			// nav.Location{
			// 	Lat: 5.5,
			// 	Lon: 5.5,
			// },
		},
	},
	[]entity.PlaceID{
		iteratortype.OrigLocationID,
		iteratortype.DestLocationID,
		1,
		1,
		1,
		2,
		2,
		2,
		3,
		3,
		3,
		4,
		4,
		4,
		5,
		5,
		5,
	},
	map[nodeID][]*edge{
		0: {
			// orig -> station 1
			{edgeID{0, 2}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{0, 3}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{0, 4}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			// orig -> station 2
			{edgeID{0, 5}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{0, 6}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{0, 7}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			// orig -> station 3
			{edgeID{0, 8}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{0, 9}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{0, 10}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
		},
		2: {
			// station 1 -> station 4
			{edgeID{2, 11}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{2, 12}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{2, 13}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			// station 1 -> station 5
			{edgeID{2, 14}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{2, 15}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{2, 16}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
		},
		3: {
			// station 1 -> station 4
			{edgeID{3, 11}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{3, 12}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{3, 13}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			// station 1 -> station 5
			{edgeID{3, 14}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{3, 15}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{3, 16}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
		},
		4: {
			// station 1 -> station 4
			{edgeID{4, 11}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{4, 12}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{4, 13}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			// station 1 -> station 5
			{edgeID{4, 14}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{4, 15}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
			{edgeID{4, 16}, &entity.Weight{Duration: 34.4, Distance: 34.4}},
		},
		5: {
			// station 2 -> station 4
			{edgeID{5, 11}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{5, 12}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{5, 13}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			// station 2 -> station 5
			{edgeID{5, 14}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{5, 15}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{5, 16}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
		},
		6: {
			// station 2 -> station 4
			{edgeID{6, 11}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{6, 12}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{6, 13}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			// station 2 -> station 5
			{edgeID{6, 14}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{6, 15}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{6, 16}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
		},
		7: {
			// station 2 -> station 4
			{edgeID{7, 11}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{7, 12}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			{edgeID{7, 13}, &entity.Weight{Duration: 11.1, Distance: 11.1}},
			// station 2 -> station 5
			{edgeID{7, 14}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{7, 15}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
			{edgeID{7, 16}, &entity.Weight{Duration: 14.4, Distance: 14.4}},
		},
		8: {
			// station 3 -> station 4
			{edgeID{8, 11}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{8, 12}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{8, 13}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			// station 3 -> station 5
			{edgeID{8, 14}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{8, 15}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{8, 16}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
		},
		9: {
			// station 3 -> station 4
			{edgeID{9, 11}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{9, 12}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{9, 13}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			// station 3 -> station 5
			{edgeID{9, 14}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{9, 15}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{9, 16}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
		},
		10: {
			// station 3 -> station 4
			{edgeID{10, 11}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{10, 12}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			{edgeID{10, 13}, &entity.Weight{Duration: 22.2, Distance: 22.2}},
			// station 3 -> station 5
			{edgeID{10, 14}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{10, 15}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
			{edgeID{10, 16}, &entity.Weight{Duration: 15.5, Distance: 15.5}},
		},
		11: {
			// station 4 -> end
			{edgeID{11, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{11, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{11, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
		},
		12: {
			// station 4 -> end
			{edgeID{12, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{12, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{12, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
		},
		13: {
			// station 4 -> end
			{edgeID{13, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{13, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
			{edgeID{13, 1}, &entity.Weight{Duration: 44.4, Distance: 44.4}},
		},
		14: {
			// station 5 -> end
			{edgeID{14, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{14, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{14, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
		},
		15: {
			// station 5 -> end
			{edgeID{15, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{15, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{15, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
		},
		16: {
			// station 5 -> end
			{edgeID{16, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{16, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
			{edgeID{16, 1}, &entity.Weight{Duration: 33.3, Distance: 33.3}},
		},
	},
	chargingstrategy.NewSimpleChargingStrategy(50.0),
}
