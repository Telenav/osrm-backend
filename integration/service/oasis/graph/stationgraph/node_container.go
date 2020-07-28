package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
)

type logicNodeIdentifier2NodePtr map[logicNodeIdentifier]*node
type nodeID2NodePtr []*node
type nodeID2PlaceID []entity.PlaceID

type nodeContainer struct {
	logicNode2NodePtr logicNodeIdentifier2NodePtr
	id2NodePtr        nodeID2NodePtr
	id2StationID      nodeID2PlaceID
	counter           int
}

func newNodeContainer() *nodeContainer {
	return &nodeContainer{
		logicNode2NodePtr: make(logicNodeIdentifier2NodePtr, 15000),
		id2NodePtr:        make(nodeID2NodePtr, 0, 15000),
		id2StationID:      make(nodeID2PlaceID, 0, 10000),
		counter:           0,
	}
}

func (nc *nodeContainer) addNode(placeID entity.PlaceID, targetState chargingstrategy.State) *node {
	key := logicNodeIdentifier{placeID, targetState}

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
		}
		nc.logicNode2NodePtr[key] = n
		nc.id2NodePtr = append(nc.id2NodePtr, n)
		nc.id2StationID = append(nc.id2StationID, placeID)
		nc.counter++

		return n
	}
}

func (nc *nodeContainer) getNode(id nodeID) *node {
	if nc.isNodeVisited(id) {
		return nc.id2NodePtr[id]
	}

	return nil
}

func (nc *nodeContainer) isNodeVisited(id nodeID) bool {
	if (int)(id) < len(nc.id2NodePtr) && nc.id2NodePtr[id] != nil {
		return true
	}
	return false
}

func (nc *nodeContainer) nodeID2PlaceID(id nodeID) entity.PlaceID {
	if (int)(id) < len(nc.id2NodePtr) {
		return nc.id2StationID[id]
	}
	return stationfindertype.InvalidPlaceID
}

type logicNodeIdentifier struct {
	placeID     entity.PlaceID
	targetState chargingstrategy.State
}
