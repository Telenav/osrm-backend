// Package waysnodes is designed to map between wayIDs and nodeIDs, e.g.,
//   `wayID->nodeIDs`
//   `fromNodeID,toNodeID->wayID(directed)`, etc.
package waysnodes

// NodesQuerier is the interface that wraps the QueryNodes method.
type NodesQuerier interface {

	// QueryNodes queries nodeIDs by wayID.
	// wayID: positive means travel forward following the nodes sequence, negative means backward
	// returned nodeIDs: the sequence will be reversed if wayID is backward.
	QueryNodes(wayID int64) (nodeIDs []int64, err error)
}

// WayQuerier is the interface that wraps the QueryWay method.
type WayQuerier interface {

	// QueryWay queries directed wayID by fromNodeID,toNodeID pair.
	// returned wayID: positive means travel forward following the fromNodeID,toNodeID sequence, negative means backward
	QueryWay(fromNodeID, toNodeID int64) (wayID int64, err error)
}

// WaysQuerier is the interface that wraps the QueryWays method.
type WaysQuerier interface {
	WayQuerier

	// QueryWays queries directed wayIDs by nodeIDs.
	// `len(wayIDs) == len(nodeIDs)-1` since each way have to be decided by traveling from one node to another.
	// returned wayIDs: positive means travel forward following the nodeIDs sequence, negative means backward
	QueryWays(nodeIDs []int64) (wayIDs []int64, err error)
}
