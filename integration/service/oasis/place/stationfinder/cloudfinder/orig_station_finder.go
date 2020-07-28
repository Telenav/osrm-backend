package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/searchhelper"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

//@todo: This number need to be adjusted based on charge station profile
const origMaxSearchCandidateNumber = 999

type origStationFinder struct {
	oasisReq *oasis.Request
	*basicFinder
}

func newOrigStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *origStationFinder {
	obj := &origStationFinder{
		oasisReq,
		newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (oFinder *origStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: oFinder.oasisReq.Coordinates[0].Lat,
			Lon: oFinder.oasisReq.Coordinates[0].Lon},
		origMaxSearchCandidateNumber,
		oFinder.oasisReq.CurrRange)

	oFinder.getNearbyChargeStations(req)
	return
}
