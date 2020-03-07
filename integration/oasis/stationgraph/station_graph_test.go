package stationgraph

import (
	"math"
	"testing"

	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
)

/*
Construct test graph of
- start connects to staion 1, station 2, station 3
- station 1 connects to station 4, station 5
- station 2 connects to station 4, station 5
- station 3 connects to station 4, station 5
- station 4 connects to end
- station 5 connects to end

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




- Each charge station will generate 3 virtual node, represent for different charge strategy:
  use different time to charge different amount of energy
- In total there should be 17 nodes in graph
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
*/
func TestConstructStationGraph(t *testing.T) {
	neighbors := [][]stationfinder.NeighborInfo{
		[]stationfinder.NeighborInfo{
			stationfinder.NeighborInfo{
				FromID: "orig_location",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 0.0,
					Lon: 0.0,
				},
				ToID: "station1",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 1.1,
					Lon: 1.1,
				},
				Cost: stationfinder.Cost{
					Duration: 22.2,
					Distance: 22.2,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "orig_location",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 0.0,
					Lon: 0.0,
				},
				ToID: "station2",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 2.2,
					Lon: 2.2,
				},
				Cost: stationfinder.Cost{
					Duration: 11.1,
					Distance: 11.1,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "orig_location",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 0.0,
					Lon: 0.0,
				},
				ToID: "station3",
				Cost: stationfinder.Cost{
					Duration: 33.3,
					Distance: 33.3,
				},
				ToLocation: stationfinder.StationCoordinate{
					Lat: 3.3,
					Lon: 3.3,
				},
			},
		},
		[]stationfinder.NeighborInfo{
			stationfinder.NeighborInfo{
				FromID: "station1",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 1.1,
					Lon: 1.1,
				},
				ToID: "station4",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 4.4,
					Lon: 4.4,
				},
				Cost: stationfinder.Cost{
					Duration: 44.4,
					Distance: 44.4,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station1",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 1.1,
					Lon: 1.1,
				},
				ToID: "station5",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 5.5,
					Lon: 5.5,
				},
				Cost: stationfinder.Cost{
					Duration: 34.3,
					Distance: 34.4,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station2",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 2.2,
					Lon: 2.2,
				},
				ToID: "station4",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 4.4,
					Lon: 4.4,
				},
				Cost: stationfinder.Cost{
					Duration: 11.1,
					Distance: 11.1,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station2",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 2.2,
					Lon: 2.2,
				},
				ToID: "station5",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 5.5,
					Lon: 5.5,
				},
				Cost: stationfinder.Cost{
					Duration: 14.4,
					Distance: 14.4,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station3",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 3.3,
					Lon: 3.3,
				},
				ToID: "station4",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 4.4,
					Lon: 4.4,
				},
				Cost: stationfinder.Cost{
					Duration: 22.2,
					Distance: 22.2,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station3",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 3.3,
					Lon: 3.3,
				},
				ToID: "station5",
				ToLocation: stationfinder.StationCoordinate{
					Lat: 5.5,
					Lon: 5.5,
				},
				Cost: stationfinder.Cost{
					Duration: 15.5,
					Distance: 15.5,
				},
			},
		},
		[]stationfinder.NeighborInfo{
			stationfinder.NeighborInfo{
				FromID: "station4",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 4.4,
					Lon: 4.4,
				},
				ToID: stationfinder.DestLocationID,
				ToLocation: stationfinder.StationCoordinate{
					Lat: 6.6,
					Lon: 6.6,
				},
				Cost: stationfinder.Cost{
					Duration: 44.4,
					Distance: 44.4,
				},
			},
			stationfinder.NeighborInfo{
				FromID: "station5",
				FromLocation: stationfinder.StationCoordinate{
					Lat: 5.5,
					Lon: 5.5,
				},
				ToID: stationfinder.DestLocationID,
				ToLocation: stationfinder.StationCoordinate{
					Lat: 6.6,
					Lon: 6.6,
				},
				Cost: stationfinder.Cost{
					Duration: 33.3,
					Distance: 33.3,
				},
			},
		},
	}

	// generate channel contains neighbors information
	// simulate real situation using different go-routine
	c := make(chan stationfinder.WeightBetweenNeighbors)
	go func() {
		for _, n := range neighbors {
			neighborsInfo := stationfinder.WeightBetweenNeighbors{
				NeighborsInfo: n,
				Err:           nil,
			}
			c <- neighborsInfo
		}
		close(c)
	}()

	currEnergyLevel := 20.0
	maxEnergyLevel := 50.0
	graph := NewStationGraph(c, currEnergyLevel, maxEnergyLevel,
		chargingstrategy.NewFakeChargingStrategyCreator(maxEnergyLevel))
	if graph == nil {
		t.Errorf("create Station graph failed, expect none-empty graph but result is empty")
	}

	testStart(t, graph, currEnergyLevel, maxEnergyLevel)
	testEnd(t, graph, currEnergyLevel, maxEnergyLevel)

	testConnectivity(t, graph, "station1", locationInfo{lat: 1.1, lon: 1.1},
		[]string{"station4", "station5"}, neighbors, currEnergyLevel, maxEnergyLevel)

	testConnectivity(t, graph, "station2", locationInfo{lat: 2.2, lon: 2.2},
		[]string{"station4", "station5"}, neighbors, currEnergyLevel, maxEnergyLevel)

	testConnectivity(t, graph, "station3", locationInfo{lat: 3.3, lon: 3.3},
		[]string{"station4", "station5"}, neighbors, currEnergyLevel, maxEnergyLevel)

	testConnectivity(t, graph, "station4", locationInfo{lat: 4.4, lon: 4.4},
		[]string{stationfinder.DestLocationID}, neighbors, currEnergyLevel, maxEnergyLevel)

	testConnectivity(t, graph, "station5", locationInfo{lat: 5.5, lon: 5.5},
		[]string{stationfinder.DestLocationID}, neighbors, currEnergyLevel, maxEnergyLevel)
}

func testStart(t *testing.T, graph *stationGraph, currEnergyLevel, maxEnergyLevel float64) {
	sn := graph.getChargeStationsNodes(stationfinder.OrigLocationID, stationfinder.StationCoordinate{}, currEnergyLevel, maxEnergyLevel)
	if len(sn) != 1 {
		t.Errorf("incorrect start node generated expect only one node but got %d", len(sn))
	}
	if graph.getStationID(sn[0].id) != stationfinder.OrigLocationID {
		t.Errorf("incorrect name for start node expect %s but got %s", stationfinder.OrigLocationID, graph.getStationID(sn[0].id))
	}
	if !floatEquals(sn[0].arrivalEnergy, currEnergyLevel) ||
		!floatEquals(sn[0].chargeEnergy, 0.0) ||
		!floatEquals(sn[0].chargeTime, 0.0) {
		t.Errorf("incorrect energy information for start node expect %v but got %v", chargeInfo{
			arrivalEnergy: currEnergyLevel,
			chargeTime:    0.0,
			chargeEnergy:  0.0,
		}, sn[0].chargeInfo)
	}

	startLocation := graph.g.getLocationInfo(sn[0].id)
	if !floatEquals(startLocation.lat, 0.0) ||
		!floatEquals(startLocation.lon, 0.0) {
		t.Errorf("incorrect location information for start node expect %v but got %v", locationInfo{
			lat: 0.0,
			lon: 0.0,
		}, startLocation)
	}

	if len(sn[0].neighbors) != 9 {
		t.Errorf("incorrect neighbors count for start node expect %d but got %d", 9, len(sn[0].neighbors))
	}

	if graph.getStationID(sn[0].neighbors[0].targetNodeID) != "station1" ||
		!floatEquals(sn[0].neighbors[0].distance, 22.2) ||
		!floatEquals(sn[0].neighbors[0].duration, 22.2) ||
		graph.getStationID(sn[0].neighbors[1].targetNodeID) != "station1" ||
		!floatEquals(sn[0].neighbors[1].distance, 22.2) ||
		!floatEquals(sn[0].neighbors[1].duration, 22.2) ||
		graph.getStationID(sn[0].neighbors[2].targetNodeID) != "station1" ||
		!floatEquals(sn[0].neighbors[2].distance, 22.2) ||
		!floatEquals(sn[0].neighbors[2].duration, 22.2) ||
		graph.getStationID(sn[0].neighbors[3].targetNodeID) != "station2" ||
		!floatEquals(sn[0].neighbors[3].distance, 11.1) ||
		!floatEquals(sn[0].neighbors[3].duration, 11.1) ||
		graph.getStationID(sn[0].neighbors[4].targetNodeID) != "station2" ||
		!floatEquals(sn[0].neighbors[4].distance, 11.1) ||
		!floatEquals(sn[0].neighbors[4].duration, 11.1) ||
		graph.getStationID(sn[0].neighbors[5].targetNodeID) != "station2" ||
		!floatEquals(sn[0].neighbors[5].distance, 11.1) ||
		!floatEquals(sn[0].neighbors[5].duration, 11.1) ||
		graph.getStationID(sn[0].neighbors[6].targetNodeID) != "station3" ||
		!floatEquals(sn[0].neighbors[6].distance, 33.3) ||
		!floatEquals(sn[0].neighbors[6].duration, 33.3) ||
		graph.getStationID(sn[0].neighbors[7].targetNodeID) != "station3" ||
		!floatEquals(sn[0].neighbors[7].distance, 33.3) ||
		!floatEquals(sn[0].neighbors[7].duration, 33.3) ||
		graph.getStationID(sn[0].neighbors[8].targetNodeID) != "station3" ||
		!floatEquals(sn[0].neighbors[8].distance, 33.3) ||
		!floatEquals(sn[0].neighbors[8].duration, 33.3) {
		t.Errorf("incorrect neighbor information generated for start node")
	}
}

func testEnd(t *testing.T, graph *stationGraph, currEnergyLevel, maxEnergyLevel float64) {
	se := graph.getChargeStationsNodes(stationfinder.DestLocationID, stationfinder.StationCoordinate{}, currEnergyLevel, maxEnergyLevel)
	if len(se) != 1 {
		t.Errorf("incorrect end node generated expect only one node but got %d", len(se))
	}
	if graph.getStationID(se[0].id) != stationfinder.DestLocationID {
		t.Errorf("incorrect name for end node expect %s but got %s", stationfinder.DestLocationID, graph.getStationID(se[0].id))
	}
	if !floatEquals(se[0].arrivalEnergy, 0.0) ||
		!floatEquals(se[0].chargeEnergy, 0.0) ||
		!floatEquals(se[0].chargeTime, 0.0) {
		t.Errorf("incorrect energy information for end node expect %v but got %v", chargeInfo{
			arrivalEnergy: 0.0,
			chargeTime:    0.0,
			chargeEnergy:  0.0,
		}, se[0].chargeInfo)
	}

	endLocation := graph.g.getLocationInfo(se[0].id)
	if !floatEquals(endLocation.lat, 6.6) ||
		!floatEquals(endLocation.lon, 6.6) {
		t.Errorf("incorrect location information for end node expect %v but got %v", locationInfo{
			lat: 6.6,
			lon: 6.6,
		}, endLocation)
	}

	if len(se[0].neighbors) != 0 {
		t.Errorf("incorrect neighbors count for end node expect %d but got %d", 0, len(se[0].neighbors))
	}
}

func testConnectivity(t *testing.T, graph *stationGraph, from string, fromLocation locationInfo,
	tos []string, mockArray [][]stationfinder.NeighborInfo, currEnergyLevel, maxEnergyLevel float64) {
	fns := graph.getChargeStationsNodes(from, stationfinder.StationCoordinate{}, 0.0, 0.0)

	if len(fns) != 3 {
		t.Errorf("incorrect node generated for %s expect 3 nodes but got %d", from, len(fns))
	}

	if !floatEquals(fns[0].arrivalEnergy, 0.0) ||
		!floatEquals(fns[0].chargeEnergy, maxEnergyLevel*0.6) ||
		!floatEquals(fns[0].chargeTime, 3600) ||
		!floatEquals(fns[1].arrivalEnergy, 0.0) ||
		!floatEquals(fns[1].chargeEnergy, maxEnergyLevel*0.8) ||
		!floatEquals(fns[1].chargeTime, 7200) ||
		!floatEquals(fns[2].arrivalEnergy, 0.0) ||
		!floatEquals(fns[2].chargeEnergy, maxEnergyLevel) ||
		!floatEquals(fns[2].chargeTime, 14400) {
		t.Errorf("incorrect charge information generated for node %s", from)
	}

	for _, fromNode := range fns {
		if !floatEquals(fromNode.locationInfo.lat, fromLocation.lat) ||
			!floatEquals(fromNode.locationInfo.lon, fromLocation.lon) {
			t.Errorf("incorrect location information generated for node %s expect %+v got %+v",
				from, fromLocation, fromNode.locationInfo)
		}
	}

	index := 0
	for _, to := range tos {
		tns := graph.getChargeStationsNodes(to, stationfinder.StationCoordinate{}, 0.0, 0.0)

		expectDuration := math.MaxFloat64
		expectDistance := math.MaxFloat64
		for _, neighborsInfo := range mockArray {
			for _, neighborInfo := range neighborsInfo {
				if neighborInfo.FromID == from && neighborInfo.ToID == to {
					expectDuration = neighborInfo.Duration
					expectDistance = neighborInfo.Distance
					break
				}
			}
		}
		if expectDuration == math.MaxFloat64 ||
			expectDistance == math.MaxFloat64 {
			t.Error("incorrect name string passed into testConnectivity")
		}

		for _, fromNode := range fns {
			for i, toNode := range tns {
				if fromNode.neighbors[index+i].targetNodeID != toNode.id ||
					fromNode.neighbors[index+i].distance != expectDistance ||
					fromNode.neighbors[index+i].duration != expectDuration {
					t.Errorf("incorrect connectivity generated between %s and %s", from, to)
				}
			}
		}

		index += len(tns)
	}
}

var epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
}
