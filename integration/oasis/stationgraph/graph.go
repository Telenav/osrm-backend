package stationgraph

import (
	"github.com/golang/glog"
)

type graph struct {
	nodes       []*node
	startNodeID nodeID
	endNodeID   nodeID
}

func (g *graph) dijkstra() []nodeID {
	m := newVisitedMap()

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
		for _, edge := range node.edges {
			if g.nodes[n].isLocationReachable(edge.distance) {
				// @todo: charge time isn't consider here, TBD
				if m.add(edge.targetNodeID, n, edge.distance, edge.duration) {
					g.nodes[edge.targetNodeID].updateArrivalEnergy(node, edge.distance)
				}
			}
		}
	}
}
