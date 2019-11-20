package wayidsmap

import "github.com/Telenav/osrm-backend/integration/graph"

// Way2Nodes defines interface to get nodeIDs from wayID.
type Way2Nodes interface {
	WayID2NodeIDs(int64) []int64
}

// Way2Edges defines interface to get Edges from wayID.
type Way2Edges interface {
	Way2Nodes

	WayID2Edges(int64) []graph.Edge
}
