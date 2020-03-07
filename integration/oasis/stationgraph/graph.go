package stationgraph

import (
	"fmt"

	"github.com/Telenav/osrm-backend/integration/oasis/chargingstrategy"
	"github.com/golang/glog"
)

type graph struct {
	nodes       []*node
	startNodeID nodeID
	endNodeID   nodeID
	strategy    chargingstrategy.ChargingStrategyCreator
}

func (g *graph) dijkstra() []nodeID {
	m := newQueryHeap()

	// init
	m.add(g.startNodeID, invalidNodeID, 0, 0)

	for {
		n := m.next()

		// stop condition
		if n == invalidNodeID {
			glog.Warning("PriorityQueue is empty before solution is found.")
			return nil
		}
		if n == g.endNodeID {
			return m.retrieve(n)
		}

		// relax
		node := g.nodes[n]
		for _, neighbor := range node.neighbors {
			fmt.Printf("+++before is reachable  %v %+v\n", g.nodes[n], neighbor)
			if g.nodes[n].isLocationReachable(neighbor.distance) {
				fmt.Printf("+++node is reachable  %v %+v\n", n, neighbor)
				chargeTime := g.nodes[neighbor.targetNodeID].calcChargeTime(node, neighbor.distance, g.strategy)
				if m.add(neighbor.targetNodeID, n, neighbor.distance, neighbor.duration+chargeTime) {
					g.nodes[neighbor.targetNodeID].updateArrivalEnergy(node, neighbor.distance)
					g.nodes[neighbor.targetNodeID].updateChargingTime(chargeTime)
				}
			}
		}
	}
}

func (g *graph) accumulateDistanceAndDuration(from nodeID, to nodeID, distance, duration *float64) {
	if from < 0 || int(from) >= len(g.nodes) {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", from)
	}

	if to < 0 || int(to) >= len(g.nodes) {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", to)
	}

	fromNode := g.nodes[from]
	for _, neighbor := range fromNode.neighbors {
		if neighbor.targetNodeID == to {
			*distance += neighbor.distance
			*duration += neighbor.duration + g.nodes[to].chargeTime
			return
		}
	}

	glog.Errorf("Passing un-connect fromNodeID and toNodeID into accumulateDistanceAndDuration.\n")
}

func (g *graph) getChargeInfo(n nodeID) chargeInfo {
	if n < 0 || int(n) >= len(g.nodes) {
		glog.Fatalf("While calling getChargeInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.nodes[n].chargeInfo
}

func (g *graph) getLocationInfo(n nodeID) locationInfo {
	if n < 0 || int(n) >= len(g.nodes) {
		glog.Fatalf("While calling getLocationInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.nodes[n].locationInfo
}
