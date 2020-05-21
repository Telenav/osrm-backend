package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/solution"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type stationGraph struct {
	g       Graph
	querier connectivitymap.Querier
}

// NewStationGraph creates station graph from channel
func NewStationGraph(currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.Strategy, querier connectivitymap.Querier) *stationGraph {
	sg := &stationGraph{
		g:       NewNodeGraph(strategy, querier),
		querier: querier,
	}

	if !sg.setStartAndEndForGraph(currEnergyLevel, maxEnergyLevel) {
		return nil
	}

	return sg
}

func (sg *stationGraph) setStartAndEndForGraph(currEnergyLevel, maxEnergyLevel float64) bool {
	startLocation := sg.querier.GetLocation(stationfindertype.OrigLocationID.String())
	if startLocation == nil {
		glog.Errorf("Failed to find %#v from Querier GetLocation()\n", stationfindertype.OrigLocationID.String())
		return false
	}

	endLocation := sg.querier.GetLocation(stationfindertype.DestLocationID.String())
	if startLocation == nil {
		glog.Errorf("Failed to find %#v from Querier GetLocation()\n", stationfindertype.DestLocationID.String())
		return false
	}

	sg.g = sg.g.SetStart(stationfindertype.OrigLocationID.String(),
		chargingstrategy.State{
			Energy: currEnergyLevel,
		},
		startLocation)

	sg.g = sg.g.SetEnd(stationfindertype.DestLocationID.String(),
		chargingstrategy.State{},
		endLocation)

	return true
}

// GenerateChargeSolutions creates charge solutions for staion graph
func (sg *stationGraph) GenerateChargeSolutions() []*solution.Solution {
	stationNodes := dijkstra(sg.g, sg.g.StartNodeID(), sg.g.EndNodeID())

	// to be removed
	//nodeGraph := sg.g.(*nodeGraph)
	//glog.Infof("+++ len(nodeGraph.adjacentList) = %v, len(nodeGraph.edgeMetric) = %v\n", len(nodeGraph.adjacentList), len(nodeGraph.edgeMetric))

	if nil == stationNodes {
		glog.Warning("Failed to generate charge stations for stationGraph.\n")
		return nil
	}

	return sg.generateSolutionsBasedOnStationCandidates(stationNodes)
}

func (sg *stationGraph) generateSolutionsBasedOnStationCandidates(stationNodes []nodeID) []*solution.Solution {
	var result []*solution.Solution

	sol := &solution.Solution{}
	sol.ChargeStations = make([]*solution.ChargeStation, 0)
	var totalDistance, totalDuration float64

	// accumulate information: start node -> first charge station
	startNodeID := sg.g.StartNodeID()
	accumulateDistanceAndDuration(sg.g, startNodeID, stationNodes[0], &totalDistance, &totalDuration)

	// accumulate information: first charge station -> second charge station -> ... -> end node
	for i := 0; i < len(stationNodes); i++ {
		if i != len(stationNodes)-1 {
			accumulateDistanceAndDuration(sg.g, stationNodes[i], stationNodes[i+1], &totalDistance, &totalDuration)
		} else {
			endNodeID := sg.g.EndNodeID()
			accumulateDistanceAndDuration(sg.g, stationNodes[i], endNodeID, &totalDistance, &totalDuration)
		}

		// construct station information
		station := &solution.ChargeStation{}
		station.ArrivalEnergy = getChargeInfo(sg.g, stationNodes[i]).arrivalEnergy
		station.ChargeRange = getChargeInfo(sg.g, stationNodes[i]).targetState.Energy
		station.ChargeTime = getChargeInfo(sg.g, stationNodes[i]).chargeTime
		station.Location = nav.Location{
			Lat: sg.getLocationInfo(sg.g, stationNodes[i]).Lat,
			Lon: sg.getLocationInfo(sg.g, stationNodes[i]).Lon,
		}
		station.StationID = sg.g.StationID(stationNodes[i])

		sol.ChargeStations = append(sol.ChargeStations, station)
	}

	sol.Distance = totalDistance
	sol.Duration = totalDuration
	sol.RemainingRage = calcRemaningRange(sg.g, stationNodes[len(stationNodes)-1])

	result = append(result, sol)
	return result
}

func accumulateDistanceAndDuration(g Graph, from nodeID, to nodeID, distance, duration *float64) {
	if g.Node(from) == nil {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", from)
	}

	if g.Node(to) == nil {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", to)
	}

	if g.Edge(from, to) == nil {
		glog.Errorf("Passing un-connect fromNodeID %#v and toNodeID %#v into accumulateDistanceAndDuration.\n", from, to)
	}

	*distance += g.Edge(from, to).distance
	*duration += g.Edge(from, to).duration + g.Node(to).chargeTime

}

// There will be different remaning range for destination node for multiple solutions
// For example:
// Solution 1: via station 111 -> station 222 -> Destination with remaning range 123
// Solution 2: via station 333 -> Destination with remaning range 45
// We decide to re-calculate remaning range instead of create multiple virtual node for destination
func calcRemaningRange(g Graph, lastStation nodeID) float64 {
	var distance, duration float64
	accumulateDistanceAndDuration(g, lastStation, g.EndNodeID(), &distance, &duration)
	return calculateArrivalEnergy(g.Node(lastStation), distance)
}

func getChargeInfo(g Graph, n nodeID) chargeInfo {
	if g.Node(n) == nil {
		glog.Fatalf("While calling getChargeInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.Node(n).chargeInfo
}

func (sg *stationGraph) getLocationInfo(g Graph, n nodeID) *nav.Location {
	if g.Node(n) == nil {
		glog.Fatalf("While calling getLocationInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return sg.querier.GetLocation(sg.g.StationID(g.Node(n).id))
}

func (sg *stationGraph) isStart(id string) bool {
	return id == stationfindertype.OrigLocationID.String()
}

func (sg *stationGraph) isEnd(id string) bool {
	return id == stationfindertype.DestLocationID.String()
}
