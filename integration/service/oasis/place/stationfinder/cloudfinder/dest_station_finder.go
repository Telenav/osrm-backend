package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/searchhelper"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

//@todo: This number need to be adjusted based on charge station profile
const destMaxSearchCandidateNumber = 999

type destStationFinder struct {
	oasisReq *oasis.Request
	*basicFinder
}

func newDestStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destStationFinder {
	obj := &destStationFinder{
		oasisReq,
		newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (dFinder *destStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: dFinder.oasisReq.Coordinates[1].Lat,
			Lon: dFinder.oasisReq.Coordinates[1].Lon},
		destMaxSearchCandidateNumber,
		dFinder.oasisReq.MaxRange-dFinder.oasisReq.SafeLevel)
	dFinder.getNearbyChargeStations(req)
}
