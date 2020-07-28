package stationgraph

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

type place2placeID struct {
	from entity.PlaceID
	to   entity.PlaceID
}
type edgeID2EdgeData map[place2placeID]*entity.Weight

func newEdgeID2EdgeData() edgeID2EdgeData {
	return make(edgeID2EdgeData, 5000000)
}

func (edge2Weight edgeID2EdgeData) get(id place2placeID) *entity.Weight {
	return edge2Weight[id]
}

func (edge2Weight edgeID2EdgeData) add(id place2placeID, weight *entity.Weight) {
	edge2Weight[id] = weight
}
