package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/golang/glog"
)

type edge struct {
	distance float64
	duration float64
}

type graph struct {
	nodes       []*node
	startNodeID nodeID
	endNodeID   nodeID
	strategy    chargingstrategy.Strategy
}

func (g *graph) Node(id nodeID) *node {
	// @todo: safety check
	return g.nodes[id]
}

func (g *graph) AdjacentList(id nodeID) []nodeID {
	// @todo: safety check
	ids := make([]nodeID, 0, len(g.nodes[id].neighbors))
	for _, neighbor := range g.nodes[id].neighbors {
		ids = append(ids, neighbor.targetNodeID)
	}
	return ids
}

func (g *graph) Edge(from, to nodeID) *edge {
	// @todo: safety check
	for _, neighbor := range g.nodes[from].neighbors {
		if neighbor.targetNodeID == to {
			return &edge{
				distance: neighbor.distance,
				duration: neighbor.duration,
			}
		}
	}
	return nil
}

func (g *graph) StartNodeID() nodeID {
	return g.startNodeID
}

func (g *graph) EndNodeID() nodeID {
	return g.endNodeID
}

func (g *graph) ChargeStrategy() chargingstrategy.Strategy {
	return g.strategy
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
