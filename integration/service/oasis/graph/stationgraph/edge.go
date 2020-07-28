package stationgraph

import "github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"

type edgeID struct {
	fromNodeID nodeID
	toNodeID   nodeID
}

type edge struct {
	edgeId     edgeID
	edgeMetric *entity.Weight
}
