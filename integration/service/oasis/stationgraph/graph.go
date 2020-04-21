package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/golang/glog"
)

// NearByStationQuery(center Location, distanceLimit float64, limitCount int) []*RankedPointInfo

type graph struct {
	nodeContainer *nodeContainer
	adjacentList  map[nodeID][]nodeID
	edgeData      map[edgeID]*edge
	startNodeID   nodeID
	endNodeID     nodeID
	strategy      chargingstrategy.Strategy
	querier       connectivitymap.Querier
}

// NewGraph creates new graph which implements IGraph
func NewGraph(strategy chargingstrategy.Strategy, query connectivitymap.Querier) IGraph {
	return &graph{
		nodeContainer: newNodeContainer(),
		adjacentList:  make(map[nodeID][]nodeID),
		edgeData:      make(map[edgeID]*edge),
		startNodeID:   invalidNodeID,
		endNodeID:     invalidNodeID,
		strategy:      strategy,
		querier:       query,
	}
}

func (g *graph) Node(id nodeID) *node {
	return g.nodeContainer.getNode(id)
}

func (g *graph) AdjacentNodes(id nodeID) []nodeID {
	if !g.nodeContainer.isNodeVisited(id) {
		glog.Errorf("While calling AdjacentNodes with un-added nodeID %#v, check your algorithm.\n", id)
		return nil
	}

	if adjList, ok := g.adjacentList[id]; ok {
		return adjList
	} else {
		adjList = g.buildAdjacentList(id)
		g.adjacentList[id] = adjList
		return adjList
	}

}

func (g *graph) Edge(from, to nodeID) *edge {
	edgeID := edgeID{
		fromNodeID: from,
		toNodeID:   to,
	}

	return g.edgeData[edgeID]
}

// SetStart generates start node for the graph
func (g *graph) SetStart(stationID string, targetState chargingstrategy.State, location locationInfo) IGraph {
	n := g.nodeContainer.addNode(stationID, targetState, location)
	g.startNodeID = n.id
	return g
}

// SetEnd generates end node for the graph
func (g *graph) SetEnd(stationID string, targetState chargingstrategy.State, location locationInfo) IGraph {
	n := g.nodeContainer.addNode(stationID, targetState, location)
	g.endNodeID = n.id
	return g
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

func (g *graph) StationID(id nodeID) string {
	return g.nodeContainer.stationID(id)
}

func (g *graph) getPhysicalAdjacentNodes(id nodeID) []*connectivitymap.QueryResult {
	stationID := g.nodeContainer.stationID(id)
	if stationID == invalidStationID {
		glog.Errorf("Query getPhysicalAdjacentNodes with invalid node %#v and result %#v\n", id, invalidStationID)
		return nil
	}
	return g.querier.NearByStationQuery(stationID)
}

func (g *graph) createLogicalNodes(from nodeID, toStationID string, toLocation *nav.Location, distance, duration float64) []*node {
	results := make([]*node, 0, 10)

	for _, state := range g.strategy.CreateChargingStates() {
		n := g.nodeContainer.addNode(toStationID, state, locationInfo{toLocation.Lat, toLocation.Lon})
		results = append(results, n)

		edgeID := edgeID{
			fromNodeID: from,
			toNodeID:   1,
		}
		g.edgeData[edgeID] = &edge{
			distance: distance,
			duration: duration,
		}
	}
	return results
}

func (g *graph) buildAdjacentList(id nodeID) []nodeID {
	adjacentNodeIDs := make([]nodeID, 0, 500)

	physicalNodes := g.getPhysicalAdjacentNodes(id)
	if physicalNodes == nil {
		glog.Errorf("Failed to build buildAdjacentList\n")
		return nil
	}

	for _, physicalNode := range physicalNodes {
		nodes := g.createLogicalNodes(id, physicalNode.StationID, physicalNode.StationLocation,
			physicalNode.Distance, physicalNode.Duration)

		for _, node := range nodes {
			adjacentNodeIDs = append(adjacentNodeIDs, node.id)
		}
	}

	return adjacentNodeIDs
}
