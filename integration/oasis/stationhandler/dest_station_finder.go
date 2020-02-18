package stationhandler

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	searchcoordinate "github.com/Telenav/osrm-backend/integration/pkg/api/search/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

//@todo: This number need to be adjusted based on charge station profile
const destMaxSearchCandidateNumber int = 999

type destStationFinder struct {
	osrmConnector     *osrmconnector.OSRMConnector
	tnSearchConnector *searchconnector.TNSearchConnector
	oasisReq          *oasis.Request
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
	bf                *basicFinder
}

func NewDestStationFinder(oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destStationFinder {
	obj := &destStationFinder{
		osrmConnector:     oc,
		tnSearchConnector: sc,
		oasisReq:          oasisReq,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
		bf:                &basicFinder{},
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

	respC := sf.tnSearchConnector.ChargeStationSearch(req)
	resp := <-respC
	if resp.Err != nil {
		glog.Warningf("Search failed during prepare orig search for url: %s", req.RequestURI())
		return
	}

	sf.searchRespLock.Lock()
	sf.searchResp = resp.Resp
	sf.searchRespLock.Unlock()
	return
}

func (sf *destStationFinder) iterateNearbyStations() <-chan chargeStationInfo {
	return sf.bf.iterateNearbyStations(sf.searchResp.Results, sf.searchRespLock)
}
