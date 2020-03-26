package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/oasis/solution"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/golang/glog"
)

type stationGraph struct {
	g               *graph
	stationID2Nodes map[string][]*node
	num2StationID   map[uint32]string // from number to original stationID
	// stationID is converted to numbers(0, 1, 2 ...) based on visit sequence

	stationsCount uint32
	strategy      chargingstrategy.Strategy
}

// NewStationGraph creates station graph from channel
func NewStationGraph(c chan stationfinder.WeightBetweenNeighbors, currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.Strategy) *stationGraph {
	sg := &stationGraph{
		g: &graph{
			startNodeID: invalidNodeID,
			endNodeID:   invalidNodeID,
			strategy:    strategy,
		},
		stationID2Nodes: make(map[string][]*node),
		num2StationID:   make(map[uint32]string),
		stationsCount:   0,
		strategy:        strategy,
	}

	for item := range c {
		if item.Err != nil {
			glog.Errorf("Met error during constructing stationgraph, error = %v", item.Err)
			return nil
		}

		for _, neighborInfo := range item.NeighborsInfo {
			sg.buildNeighborInfoBetweenNodes(neighborInfo, currEnergyLevel, maxEnergyLevel)
		}
	}

	return sg.constructGraph()
}

// GenerateChargeSolutions creates creates charge solutions for staion graph
func (sg *stationGraph) GenerateChargeSolutions() []*solution.Solution {
	stationNodes := sg.g.dijkstra()
	if nil == stationNodes {
		glog.Warning("Failed to generate charge stations for stationGraph.\n")
		return nil
	}

	var result []*solution.Solution

	sol := &solution.Solution{}
	sol.ChargeStations = make([]*solution.ChargeStation, 0)
	var totalDistance, totalDuration float64

	// accumulate information: start node -> first charge station
	startNodeID := sg.stationID2Nodes[stationfinder.OrigLocationID][0].id
	sg.g.accumulateDistanceAndDuration(startNodeID, stationNodes[0], &totalDistance, &totalDuration)

	// accumulate information: first charge station -> second charge station -> ... -> end node
	for i := 0; i < len(stationNodes); i++ {
		if i != len(stationNodes)-1 {
			sg.g.accumulateDistanceAndDuration(stationNodes[i], stationNodes[i+1], &totalDistance, &totalDuration)
		} else {
			endNodeID := sg.stationID2Nodes[stationfinder.DestLocationID][0].id
			sg.g.accumulateDistanceAndDuration(stationNodes[i], endNodeID, &totalDistance, &totalDuration)
		}

		// construct station information
		station := &solution.ChargeStation{}
		station.ArrivalEnergy = sg.g.getChargeInfo(stationNodes[i]).arrivalEnergy
		station.ChargeRange = sg.g.getChargeInfo(stationNodes[i]).targetState.Energy
		station.ChargeTime = sg.g.getChargeInfo(stationNodes[i]).chargeTime
		station.Location = nav.Location{
			Lat: sg.g.getLocationInfo(stationNodes[i]).lat,
			Lon: sg.g.getLocationInfo(stationNodes[i]).lon,
		}
		station.StationID = sg.num2StationID[uint32(stationNodes[i])]

		sol.ChargeStations = append(sol.ChargeStations, station)
	}

	sol.Distance = totalDistance
	sol.Duration = totalDuration
	sol.RemainingRage = sg.stationID2Nodes[stationfinder.DestLocationID][0].arrivalEnergy

	result = append(result, sol)
	return result
}

func (sg *stationGraph) buildNeighborInfoBetweenNodes(neighborInfo stationfinder.NeighborInfo, currEnergyLevel, maxEnergyLevel float64) {
	for _, fromNode := range sg.getChargeStationsNodes(neighborInfo.FromID, neighborInfo.FromLocation, currEnergyLevel, maxEnergyLevel) {
		for _, toNode := range sg.getChargeStationsNodes(neighborInfo.ToID, neighborInfo.ToLocation, currEnergyLevel, maxEnergyLevel) {
			fromNode.neighbors = append(fromNode.neighbors, &neighbor{
				targetNodeID: toNode.id,
				distance:     neighborInfo.Distance,
				duration:     neighborInfo.Duration,
			})
		}
	}
}

func (sg *stationGraph) getChargeStationsNodes(id string, location nav.Location, currEnergyLevel, maxEnergyLevel float64) []*node {
	if _, ok := sg.stationID2Nodes[id]; !ok {
		if sg.isStart(id) {
			sg.constructStartNode(id, location, currEnergyLevel)
		} else if sg.isEnd(id) {
			sg.constructEndNode(id, location)
		} else {
			var nodes []*node
			for _, state := range sg.strategy.CreateChargingStates() {
				n := &node{
					id: nodeID(sg.stationsCount),
					chargeInfo: chargeInfo{
						targetState: state,
					},
					locationInfo: locationInfo{
						lat: location.Lat,
						lon: location.Lon,
					},
				}
				nodes = append(nodes, n)
				sg.num2StationID[sg.stationsCount] = id
				sg.stationsCount += 1
			}

			sg.stationID2Nodes[id] = nodes
		}
	}
	return sg.stationID2Nodes[id]
}

func (sg *stationGraph) isStart(id string) bool {
	return id == stationfinder.OrigLocationID
}

func (sg *stationGraph) isEnd(id string) bool {
	return id == stationfinder.DestLocationID
}

func (sg *stationGraph) getStationID(id nodeID) string {
	return sg.num2StationID[uint32(id)]
}

func (sg *stationGraph) constructStartNode(id string, location nav.Location, currEnergyLevel float64) {

	n := &node{
		id: nodeID(sg.stationsCount),
		chargeInfo: chargeInfo{
			arrivalEnergy: currEnergyLevel,
			chargeTime:    0.0,
			targetState: chargingstrategy.State{
				Energy: currEnergyLevel,
			},
		},
		locationInfo: locationInfo{
			lat: location.Lat,
			lon: location.Lon,
		},
	}
	sg.stationID2Nodes[id] = []*node{n}
	sg.num2StationID[sg.stationsCount] = id
	sg.stationsCount += 1
}

func (sg *stationGraph) constructEndNode(id string, location nav.Location) {

	n := &node{
		id: nodeID(sg.stationsCount),
		locationInfo: locationInfo{
			lat: location.Lat,
			lon: location.Lon,
		},
	}
	sg.stationID2Nodes[id] = []*node{n}
	sg.num2StationID[sg.stationsCount] = id
	sg.stationsCount += 1
}

func (sg *stationGraph) constructGraph() *stationGraph {
	sg.g.nodes = make([]*node, int(sg.stationsCount))

	for k, v := range sg.stationID2Nodes {
		if sg.isStart(k) {
			sg.g.startNodeID = v[0].id
		}

		if sg.isEnd(k) {
			sg.g.endNodeID = v[0].id
		}

		for _, n := range v {
			sg.g.nodes[n.id] = n
		}
	}

	if sg.g.startNodeID == invalidNodeID {
		glog.Error("Invalid nodeid generated for start node.\n")
		return nil
	} else if sg.g.endNodeID == invalidNodeID {
		glog.Error("Invalid nodeid generated for start node.\n")
		return nil
	} else if len(sg.g.nodes) != int(sg.stationsCount) {
		glog.Errorf("Invalid nodes generated, len(sg.g.nodes) is %d while sg.stationsCount is %d.\n", len(sg.g.nodes), sg.stationsCount)
		return nil
	}

	return sg
}
