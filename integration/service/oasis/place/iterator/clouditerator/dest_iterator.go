package clouditerator

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/searchhelper"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

//@todo: This number need to be adjusted based on charge station profile
const destMaxSearchCandidateNumber = 999

type destIterator struct {
	oasisReq *oasis.Request
	*basicIterator
}

func newDestIterator(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destIterator {
	obj := &destIterator{
		oasisReq,
		newBasicIterator(sc),
	}
	obj.prepare()
	return obj
}

func (dFinder *destIterator) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: dFinder.oasisReq.Coordinates[1].Lat,
			Lon: dFinder.oasisReq.Coordinates[1].Lon},
		destMaxSearchCandidateNumber,
		dFinder.oasisReq.MaxRange-dFinder.oasisReq.SafeLevel)
	dFinder.getNearbyChargeStations(req)
}
