package stationgraph

import "github.com/golang/glog"

// dijkstra accepts Graph and returns all node ids in the shortest path, except start node and end node
func dijkstra(g Graph, start, end nodeID) []nodeID {
	m := newQueryHeap()

	// init
	m.add(start, invalidNodeID, 0, 0)

	for {
		currID := m.next()

		// stop condition
		if currID == invalidNodeID {
			glog.Warning("PriorityQueue is empty before solution is found.")
			return nil
		}
		if currID == end {
			// to be removed
			//glog.Infof("+++  len(queryHeap.m) = %v \n", len(m.m))

			return m.retrieve(currID)
		}

		// relax
		node := g.Node(currID)
		for _, targetID := range g.AdjacentNodes(currID) {
			if g.Edge(currID, targetID) == nil {
				glog.Errorf("No connectivity between %+v and %+v which is unexpected, check your logic.\n", currID, targetID)
			}

			len := g.Edge(currID, targetID).distance
			t := g.Edge(currID, targetID).duration

			// todo
			// !g.Node(currID).tooNearToCurrentNode(g.Node(targetID), len)
			if g.Node(currID).reachableByDistance(len) {
				chargeTimeNeeded := g.Node(targetID).calcChargeTime(node, len, g.ChargeStrategy())
				if m.add(targetID, currID, len, chargeTimeNeeded+t) {
					g.Node(targetID).updateArrivalEnergy(node, len)
					g.Node(targetID).updateChargingTime(chargeTimeNeeded)
				}
			}
		}
	}
}
