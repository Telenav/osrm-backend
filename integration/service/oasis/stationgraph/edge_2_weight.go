package stationgraph

import "github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"

type place2placeID struct {
	from common.PlaceID
	to   common.PlaceID
}
type edgeID2EdgeData map[place2placeID]*common.Weight

func newEdgeID2EdgeData() edgeID2EdgeData {
	return make(edgeID2EdgeData, 5000000)
}

func (edge2Weight edgeID2EdgeData) get(id place2placeID) *common.Weight {
	return edge2Weight[id]
}

func (edge2Weight edgeID2EdgeData) add(id place2placeID, weight *common.Weight) {
	edge2Weight[id] = weight
}
