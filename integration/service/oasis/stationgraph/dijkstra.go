package stationgraph

import "github.com/golang/glog"

// dijkstra accepts IGraph and returns all node ids in the shortest path, except start node and end node
func dijkstra(g IGraph) []nodeID {
	m := newQueryHeap()

	// init
	m.add(g.StartNodeID(), invalidNodeID, 0, 0)

	for {
		currID := m.next()

		// stop condition
		if currID == invalidNodeID {
			glog.Warning("PriorityQueue is empty before solution is found.")
			return nil
		}
		if currID == g.EndNodeID() {
			return m.retrieve(currID)
		}

		// relax
		node := g.Node(currID)
		for _, targetID := range g.AdjacentList(currID) {
			if g.Edge(currID, targetID) == nil {
				glog.Errorf("No connectivity between %#v and %#v which is unexpected, check your logic.\n", currID, targetID)
			}

			len := g.Edge(currID, targetID).distance
			t := g.Edge(currID, targetID).duration

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
