package stationfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
)

//@todo: This number need to be adjusted based on charge station profile
const destMaxSearchCandidateNumber = 999

type destStationFinder struct {
	tnSearchConnector *searchconnector.TNSearchConnector
	oasisReq          *oasis.Request
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
	bf                *basicFinder
}

func NewDestStationFinder(sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destStationFinder {
	obj := &destStationFinder{
		oasisReq: oasisReq,
		bf:       newBasicFinder(sc),
	}
	obj.prepare()
	return obj
}

func (sf *destStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.oasisReq.Coordinates[1].Lat,
			Lon: sf.oasisReq.Coordinates[1].Lon},
		destMaxSearchCandidateNumber,
		sf.oasisReq.MaxRange-sf.oasisReq.SafeLevel)
	sf.bf.getNearbyChargeStations(req)

	return
}

func (sf *destStationFinder) iterateNearbyStations() <-chan ChargeStationInfo {
	return sf.bf.iterateNearbyStations()
}
