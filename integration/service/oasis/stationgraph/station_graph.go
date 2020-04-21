package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/solution"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type stationGraph struct {
	g        IGraph
	querier  connectivitymap.Querier
	strategy chargingstrategy.Strategy
}

// NewStationGraph creates station graph from channel
func NewStationGraph(currEnergyLevel, maxEnergyLevel float64, strategy chargingstrategy.Strategy, querier connectivitymap.Querier) *stationGraph {
	sg := &stationGraph{
		g:        NewGraph(strategy, querier),
		querier:  querier,
		strategy: strategy,
	}

	if !sg.setStartAndEnd(currEnergyLevel, maxEnergyLevel) {
		return nil
	}

	return sg
}

func (sg *stationGraph) setStartAndEnd(currEnergyLevel, maxEnergyLevel float64) bool {
	startLocation := sg.querier.GetLocation(stationfindertype.OrigLocationID)
	if startLocation == nil {
		glog.Errorf("Failed to find %#v from Querier's GetLocation()\n", stationfindertype.OrigLocationID)
		return false
	}

	endLocation := sg.querier.GetLocation(stationfindertype.DestLocationID)
	if startLocation == nil {
		glog.Errorf("Failed to find %#v from Querier's GetLocation()\n", stationfindertype.DestLocationID)
		return false
	}

	sg.g = sg.g.SetStart(stationfindertype.OrigLocationID,
		chargingstrategy.State{
			Energy: currEnergyLevel,
		},
		locationInfo{
			startLocation.Lat,
			startLocation.Lon})

	sg.g = sg.g.SetEnd(stationfindertype.DestLocationID,
		chargingstrategy.State{},
		locationInfo{
			endLocation.Lat,
			endLocation.Lon})

	return true
}

// GenerateChargeSolutions creates creates charge solutions for staion graph
func (sg *stationGraph) GenerateChargeSolutions() []*solution.Solution {
	stationNodes := dijkstra(sg.g, sg.g.StartNodeID(), sg.g.EndNodeID())
	if nil == stationNodes {
		glog.Warning("Failed to generate charge stations for stationGraph.\n")
		return nil
	}

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
			Lat: getLocationInfo(sg.g, stationNodes[i]).lat,
			Lon: getLocationInfo(sg.g, stationNodes[i]).lon,
		}
		station.StationID = sg.g.StationID(stationNodes[i])

		sol.ChargeStations = append(sol.ChargeStations, station)
	}

	sol.Distance = totalDistance
	sol.Duration = totalDuration
	sol.RemainingRage = getChargeInfo(sg.g, sg.g.EndNodeID()).arrivalEnergy

	result = append(result, sol)
	return result
}

func accumulateDistanceAndDuration(g IGraph, from nodeID, to nodeID, distance, duration *float64) {
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

func getChargeInfo(g IGraph, n nodeID) chargeInfo {
	if g.Node(n) == nil {
		glog.Fatalf("While calling getChargeInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.Node(n).chargeInfo
}

func getLocationInfo(g IGraph, n nodeID) locationInfo {
	if g.Node(n) == nil {
		glog.Fatalf("While calling getLocationInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.Node(n).locationInfo
}

func (sg *stationGraph) isStart(id string) bool {
	return id == stationfindertype.OrigLocationID
}

func (sg *stationGraph) isEnd(id string) bool {
	return id == stationfindertype.DestLocationID
}
