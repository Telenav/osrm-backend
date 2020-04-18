package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/golang/glog"
)

type logicNodeIdentifier2NodePtrMap map[logicNodeIdentifier]*node
type nodeID2NodePtr map[nodeID]*node
type nodeID2StationID map[nodeID]string

const invalidStationID = "InvalidStationID"

type logicNodeIdentifier struct {
	stationID   string
	targetState chargingstrategy.State
}

type nodeContainer struct {
	logicNode2NodePtr logicNodeIdentifier2NodePtrMap
	id2NodePtr        nodeID2NodePtr
	id2StationID      nodeID2StationID
	counter           int
}

func newNodeContainer() *nodeContainer {
	return &nodeContainer{
		logicNode2NodePtr: make(logicNodeIdentifier2NodePtrMap),
		id2NodePtr:        make(nodeID2NodePtr),
		counter:           0,
	}
}

func (nc *nodeContainer) addNode(stationID string, targetState chargingstrategy.State, location locationInfo) *node {
	key := logicNodeIdentifier{stationID, targetState}

	if n, ok := nc.logicNode2NodePtr[key]; ok {
		return n
	} else {
		n = &node{
			id: (nodeID(nc.counter)),
			chargeInfo: chargeInfo{
				arrivalEnergy: 0.0,
				chargeTime:    0.0,
				targetState:   targetState,
			},
			locationInfo: location,
		}
		nc.logicNode2NodePtr[key] = n
		nc.id2NodePtr[n.id] = n
		nc.id2StationID[n.id] = stationID
		nc.counter++

		return n
	}
}

func (nc *nodeContainer) getNode(id nodeID) *node {
	if n, ok := nc.id2NodePtr[id]; ok {
		return n
	} else {
		return nil
	}

}

func (nc *nodeContainer) isNodeVisited(id nodeID) bool {
	_, ok := nc.id2NodePtr[id]
	return ok
}

func (nc *nodeContainer) getStationID(id nodeID) string {
	if stationID, ok := nc.id2StationID[id]; ok {
		return stationID
	} else {
		return invalidStationID
	}
}

// NearByStationQuery(center Location, distanceLimit float64, limitCount int) []*RankedPointInfo

type graph struct {
	nodes         []*node
	nodeContainer nodeContainer
	adjacentList  map[nodeID][]nodeID
	edgeData      map[edgeID]*edge
	startNodeID   nodeID
	endNodeID     nodeID
	strategy      chargingstrategy.Strategy
	query         connectivitymap.Querier
}

func (g *graph) Node(id nodeID) *node {
	return g.nodeContainer.getNode(id)
}

func (g *graph) getPhysicalAdjacentNodes(id nodeID) []*connectivitymap.QueryResult {
	stationID := g.nodeContainer.getStationID(id)
	if stationID == invalidStationID {
		glog.Errorf("Query getPhysicalAdjacentNodes with invalid node %#v and result %#v\n", id, invalidStationID)
		return nil
	}
	return g.query.NearByStationQuery(stationID)
}

func (g *graph) createLogicalNodes(from nodeID, toStationID string, toLocation nav.Location, distance, duration float64) []*node {
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

func (g *graph) StartNodeID() nodeID {
	return g.startNodeID
}

func (g *graph) EndNodeID() nodeID {
	return g.endNodeID
}

func (g *graph) ChargeStrategy() chargingstrategy.Strategy {
	return g.strategy
}

func (g *graph) accumulateDistanceAndDuration(from nodeID, to nodeID, distance, duration *float64) {
	if from < 0 || int(from) >= len(g.nodes) {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", from)
	}

	if to < 0 || int(to) >= len(g.nodes) {
		glog.Fatalf("While calling accumulateDistanceAndDuration, incorrect nodeID passed into graph %v\n", to)
	}

	if g.Edge(from, to) == nil {
		glog.Errorf("Passing un-connect fromNodeID %#v and toNodeID %#v into accumulateDistanceAndDuration.\n", from, to)
	}

	*distance += g.Edge(from, to).distance
	*duration += g.Edge(from, to).duration + g.Node(to).chargeTime

}

func (g *graph) getChargeInfo(n nodeID) chargeInfo {
	if n < 0 || int(n) >= len(g.nodes) {
		glog.Fatalf("While calling getChargeInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.nodes[n].chargeInfo
}

func (g *graph) getLocationInfo(n nodeID) locationInfo {
	if n < 0 || int(n) >= len(g.nodes) {
		glog.Fatalf("While calling getLocationInfo, incorrect nodeID passed into graph %v\n", n)
	}

	return g.nodes[n].locationInfo
}
