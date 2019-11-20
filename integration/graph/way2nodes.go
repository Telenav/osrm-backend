package graph

// WayID2NodeIDsMapping defines interface to get nodeIDs from wayID.
type WayID2NodeIDsMapping interface {
	GetNodeIDs(int64) []int64
}

// WayID2EdgesMapping defines interface to get Edges from wayID.
type WayID2EdgesMapping interface {
	WayID2NodeIDsMapping

	GetEdges(int64) []Edge
}
