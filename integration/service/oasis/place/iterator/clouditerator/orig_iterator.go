package clouditerator

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/searchhelper"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

//@todo: This number need to be adjusted based on charge station profile
const origMaxSearchCandidateNumber = 999

type origIterator struct {
	oasisReq *oasis.Request
	*basicIterator
}

func newOrigIterator(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *origIterator {
	obj := &origIterator{
		oasisReq,
		newBasicIterator(sc),
	}
	obj.prepare()
	return obj
}

func (oFinder *origIterator) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: oFinder.oasisReq.Coordinates[0].Lat,
			Lon: oFinder.oasisReq.Coordinates[0].Lon},
		origMaxSearchCandidateNumber,
		oFinder.oasisReq.CurrRange)

	oFinder.getNearbyChargeStations(req)
	return
}
