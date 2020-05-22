package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

// Graph defines interface used for Graph
type Graph interface {

	// Node returns node object by its nodeID
	Node(id nodeID) *node

	// AdjacentNodes returns a group of node ids which connect with given node id
	// The connectivity between nodes is build during running time.
	AdjacentNodes(id nodeID) []nodeID

	// Edge returns edge information between given two nodes
	Edge(from, to nodeID) *common.Weight

	// SetStart generates start node for the graph
	SetStart(placeID common.PlaceID, targetState chargingstrategy.State, location *nav.Location) Graph

	// SetEnd generates end node for the graph
	SetEnd(placeID common.PlaceID, targetState chargingstrategy.State, location *nav.Location) Graph

	// StartNodeID returns start node's ID for given graph
	StartNodeID() nodeID

	// EndNodeID returns end node's ID for given graph
	EndNodeID() nodeID

	// ChargeStrategy returns charge strategy used for graph construction
	ChargeStrategy() chargingstrategy.Strategy

	// PlaceID returns original placeID from internal nodeID
	PlaceID(id nodeID) common.PlaceID
}
