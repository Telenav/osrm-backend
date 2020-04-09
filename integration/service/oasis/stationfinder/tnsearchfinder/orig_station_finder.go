package tnsearchfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

//@todo: This number need to be adjusted based on charge station profile
const origMaxSearchCandidateNumber = 999

type origStationFinder struct {
	oasisReq *oasis.Request
	bf       *basicFinder
}

func NewOrigStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *origStationFinder {
	obj := &origStationFinder{
		oasisReq: oasisReq,
		bf:       newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (sf *origStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.oasisReq.Coordinates[0].Lat,
			Lon: sf.oasisReq.Coordinates[0].Lon},
		origMaxSearchCandidateNumber,
		sf.oasisReq.CurrRange)

	sf.bf.getNearbyChargeStations(req)
	return
}

func (sf *origStationFinder) IterateNearbyStations() <-chan stationfindertype.ChargeStationInfo {
	return sf.bf.IterateNearbyStations()
}
