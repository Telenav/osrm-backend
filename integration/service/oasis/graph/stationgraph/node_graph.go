package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type nodeID2AdjacentNodes map[nodeID][]nodeID

type nodeGraph struct {
	nodeContainer *nodeContainer
	adjacentList  nodeID2AdjacentNodes
	edgeMetric    edgeID2EdgeData
	startNodeID   nodeID
	endNodeID     nodeID
	strategy      chargingstrategy.Strategy
	querier       place.TopoQuerier
}

// NewNodeGraph creates new node based graph which implements Graph
func NewNodeGraph(strategy chargingstrategy.Strategy, query place.TopoQuerier) Graph {
	return &nodeGraph{
		nodeContainer: newNodeContainer(),
		adjacentList:  make(nodeID2AdjacentNodes),
		edgeMetric:    newEdgeID2EdgeData(),
		startNodeID:   invalidNodeID,
		endNodeID:     invalidNodeID,
		strategy:      strategy,
		querier:       query,
	}
}

// Node returns node object by its nodeID
func (g *nodeGraph) Node(id nodeID) *node {
	return g.nodeContainer.getNode(id)
}

// AdjacentNodes returns a group of node ids which connect with given node id
// The connectivity between nodes is build during running time.
func (g *nodeGraph) AdjacentNodes(id nodeID) []nodeID {
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

// Edge returns edge information between given two nodes
func (g *nodeGraph) Edge(from, to nodeID) *entity.Weight {
	return g.edgeMetric.get(place2placeID{
		from: g.nodeContainer.nodeID2PlaceID(from),
		to:   g.nodeContainer.nodeID2PlaceID(to)})
}

// SetStart generates start node for the nodeGraph
func (g *nodeGraph) SetStart(placeID entity.PlaceID, targetState chargingstrategy.State, location *nav.Location) Graph {
	n := g.nodeContainer.addNode(placeID, targetState)
	g.startNodeID = n.id
	return g
}

// SetEnd generates end node for the nodeGraph
func (g *nodeGraph) SetEnd(placeID entity.PlaceID, targetState chargingstrategy.State, location *nav.Location) Graph {
	n := g.nodeContainer.addNode(placeID, targetState)
	g.endNodeID = n.id
	return g
}

// StartNodeID returns start node's ID for given graph
func (g *nodeGraph) StartNodeID() nodeID {
	return g.startNodeID
}

// EndNodeID returns end node's ID for given graph
func (g *nodeGraph) EndNodeID() nodeID {
	return g.endNodeID
}

// ChargeStrategy returns charge strategy used for graph construction
func (g *nodeGraph) ChargeStrategy() chargingstrategy.Strategy {
	return g.strategy
}

// PlaceID returns original placeID from internal nodeID
func (g *nodeGraph) PlaceID(id nodeID) entity.PlaceID {
	return g.nodeContainer.nodeID2PlaceID(id)
}

func (g *nodeGraph) getPhysicalAdjacentNodes(id nodeID) []*entity.RankedPlaceInfo {
	placeID := g.nodeContainer.nodeID2PlaceID(id)
	if placeID == stationfindertype.InvalidPlaceID {
		glog.Errorf("Query getPhysicalAdjacentNodes with invalid node %#v and result %#v\n", id, stationfindertype.InvalidPlaceIDStr)
		return nil
	}
	return g.querier.NearByStationQuery(placeID)
}

func (g *nodeGraph) createLogicalNodes(from nodeID, toPlaceID entity.PlaceID, toLocation *nav.Location, weight *entity.Weight) []*node {
	results := make([]*node, 0, 3)

	// if toPlaceID equals endNode, direct return since there is no need to create charge candidates for endNode
	endNodeID := g.EndNodeID()
	if toPlaceID == g.PlaceID(endNodeID) {
		results = append(results, g.Node(endNodeID))
		g.edgeMetric.add(place2placeID{
			from: g.nodeContainer.nodeID2PlaceID(from),
			to:   g.nodeContainer.nodeID2PlaceID(endNodeID)}, weight)
		return results
	}

	// creates logical nodes for each physical node, different logical node means charge for different level of energy at charge station
	for _, state := range g.strategy.CreateChargingStates() {
		n := g.nodeContainer.addNode(toPlaceID, state)
		results = append(results, n)

		g.edgeMetric.add(place2placeID{
			from: g.nodeContainer.nodeID2PlaceID(from),
			to:   g.nodeContainer.nodeID2PlaceID(n.id)}, weight)
	}
	return results
}

func (g *nodeGraph) buildAdjacentList(id nodeID) []nodeID {

	physicalNodes := g.getPhysicalAdjacentNodes(id)
	if physicalNodes == nil {
		glog.Errorf("Failed to build buildAdjacentList.\n")
		return nil
	}

	numOfPhysicalNodesNeeded := 0
	for _, physicalNode := range physicalNodes {
		// filter nodes which is un-reachable by current energy, nodes are sorted based on distance
		if !g.Node(id).reachableByDistance(physicalNode.Weight.Distance) {
			break
		}
		numOfPhysicalNodesNeeded++
	}
	adjacentNodeIDs := make([]nodeID, 0, numOfPhysicalNodesNeeded*3)

	for _, physicalNode := range physicalNodes {
		// filter nodes which is un-reachable by current energy, nodes are sorted based on distance
		if !g.Node(id).reachableByDistance(physicalNode.Weight.Distance) {
			break
		}

		nodes := g.createLogicalNodes(id, physicalNode.ID, physicalNode.Location,
			physicalNode.Weight)

		for _, node := range nodes {
			adjacentNodeIDs = append(adjacentNodeIDs, node.id)
		}
	}

	return adjacentNodeIDs
}
