package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
)

type logicNodeIdentifier2NodePtr map[logicNodeIdentifier]*node
type nodeID2NodePtr map[nodeID]*node
type nodeID2StationID map[nodeID]string

type nodeContainer struct {
	logicNode2NodePtr logicNodeIdentifier2NodePtr
	id2NodePtr        nodeID2NodePtr
	id2StationID      nodeID2StationID
	counter           int
}

func newNodeContainer() *nodeContainer {
	return &nodeContainer{
		logicNode2NodePtr: make(logicNodeIdentifier2NodePtr),
		id2NodePtr:        make(nodeID2NodePtr),
		id2StationID:      make(nodeID2StationID),
		counter:           0,
	}
}

func (nc *nodeContainer) addNode(stationID string, targetState chargingstrategy.State, location nav.Location) *node {
	key := logicNodeIdentifier{stationID, targetState}

	if n, ok := nc.logicNode2NodePtr[key]; ok {
		return n
	} else {
		n = &node{
			(nodeID(nc.counter)),
			chargeInfo{
				arrivalEnergy: 0.0,
				chargeTime:    0.0,
				targetState:   targetState,
			},
			nav.Location{
				Lat: location.Lat,
				Lon: location.Lon},
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
	}

	return nil
}

func (nc *nodeContainer) isNodeVisited(id nodeID) bool {
	_, ok := nc.id2NodePtr[id]
	return ok
}

func (nc *nodeContainer) stationID(id nodeID) string {
	if stationID, ok := nc.id2StationID[id]; ok {
		return stationID
	} else {
		return invalidStationID
	}
}

const invalidStationID = "InvalidStationID"

type logicNodeIdentifier struct {
	stationID   string
	targetState chargingstrategy.State
}
