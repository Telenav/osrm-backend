package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/oasis/solutionformat"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/golang/glog"
)

type stationGraph struct {
	g               *graph
	stationID2Nodes map[string][]*node
	num2StationID   map[uint32]string // from number to original stationID
	// stationID is converted to numbers(0, 1, 2 ...) based on visit sequence

	stationsCount uint32
	strategy      chargingstrategy.ChargingStrategyCreator
}

// NewStationGraph creates station graph from channel
func NewStationGraph(c chan stationfinder.WeightBetweenNeighbors, currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.ChargingStrategyCreator) *stationGraph {
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

func (sg *stationGraph) GenerateChargeSolutions() []*solutionformat.Solution {
	nodes := sg.g.dijkstra()
	if nil == nodes {
		glog.Warning("Failed to generate charge stations for stationGraph.\n")
		return nil
	}

	var result []*solutionformat.Solution

	solution := &solutionformat.Solution{}
	solution.ChargeStations = make([]*solutionformat.ChargeStation, 0)
	var totalDistance, totalDuration float64

	startNodeID := sg.stationID2Nodes[stationfinder.OrigLocationID][0].id
	sg.g.accumulateDistanceAndDuration(startNodeID, nodes[0], &totalDistance, &totalDuration)
	for i := 0; i < len(nodes); i++ {
		if i != len(nodes)-1 {
			sg.g.accumulateDistanceAndDuration(nodes[i], nodes[i+1], &totalDistance, &totalDuration)
		} else {
			endNodeID := sg.stationID2Nodes[stationfinder.DestLocationID][0].id
			sg.g.accumulateDistanceAndDuration(nodes[i], endNodeID, &totalDistance, &totalDuration)
		}

		station := &solutionformat.ChargeStation{}
		station.ArrivalEnergy = sg.g.getChargeInfo(nodes[i]).arrivalEnergy
		station.ChargeRange = sg.g.getChargeInfo(nodes[i]).chargeEnergy
		station.ChargeTime = sg.g.getChargeInfo(nodes[i]).chargeTime
		station.Location = solutionformat.Location{
			Lat: sg.g.getLocationInfo(nodes[i]).lat,
			Lon: sg.g.getLocationInfo(nodes[i]).lon,
		}
		station.StationID = sg.num2StationID[uint32(nodes[i])]

		solution.ChargeStations = append(solution.ChargeStations, station)

	}

	solution.Distance = totalDistance
	solution.Duration = totalDuration
	solution.RemainingRage = sg.stationID2Nodes[stationfinder.DestLocationID][0].arrivalEnergy

	result = append(result, solution)
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

func (sg *stationGraph) getChargeStationsNodes(id string, location stationfinder.StationCoordinate, currEnergyLevel, maxEnergyLevel float64) []*node {
	if _, ok := sg.stationID2Nodes[id]; !ok {
		if sg.isStart(id) {
			sg.constructStartNode(id, location, currEnergyLevel)
		} else if sg.isEnd(id) {
			sg.constructEndNode(id, location)
		} else {
			var nodes []*node
			for _, strategy := range sg.strategy.CreateChargingStrategies() {
				n := &node{
					id: nodeID(sg.stationsCount),
					chargeInfo: chargeInfo{
						chargeEnergy: strategy.ChargingEnergy,
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

func (sg *stationGraph) constructStartNode(id string, location stationfinder.StationCoordinate, currEnergyLevel float64) {

	n := &node{
		id:         nodeID(sg.stationsCount),
		chargeInfo: chargeInfo{arrivalEnergy: currEnergyLevel},
		locationInfo: locationInfo{
			lat: location.Lat,
			lon: location.Lon,
		},
	}
	sg.stationID2Nodes[id] = []*node{n}
	sg.num2StationID[sg.stationsCount] = id
	sg.stationsCount += 1
}

func (sg *stationGraph) constructEndNode(id string, location stationfinder.StationCoordinate) {

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
