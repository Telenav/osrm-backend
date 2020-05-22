package stationgraph

import "github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"

type edgeID struct {
	fromNodeID nodeID
	toNodeID   nodeID
}

type edge struct {
	edgeId     edgeID
	edgeMetric *common.Weight
}
