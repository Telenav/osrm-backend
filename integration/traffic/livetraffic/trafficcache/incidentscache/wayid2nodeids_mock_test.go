package incidentscache

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/graph"
)

type wayID2NodeIDs map[int64][]int64

func (w wayID2NodeIDs) WayID2NodeIDs(wayID int64) []int64 {
	nodeIDs, found := w[wayID]
	if found {
		return nodeIDs
	}
	return nil
}

func (w wayID2NodeIDs) WayID2Edges(wayID int64) []graph.Edge {

	absWayID := int64(math.Abs(float64(wayID)))
	nodeIDs, found := w[absWayID]
	if found {
		edges := []graph.Edge{}
		for i := range nodeIDs[:len(nodeIDs)-1] {
			edges = append(edges, graph.Edge{From: nodeIDs[i], To: nodeIDs[i+1]})
		}

		if wayID < 0 {
			return graph.ReverseEdges(edges)
		}
		return edges
	}
	return nil
}
