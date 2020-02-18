package stationhandler

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	searchcoordinate "github.com/Telenav/osrm-backend/integration/pkg/api/search/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

type destStationFinder struct {
	osrmConnector     *osrmconnector.OSRMConnector
	tnSearchConnector *searchconnector.TNSearchConnector
	oasisReq          *oasis.Request
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
}

func NewDestStationFinder(oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector, oasisReq *oasis.Request) *destStationFinder {
	obj := &destStationFinder{
		osrmConnector:     oc,
		tnSearchConnector: sc,
		oasisReq:          oasisReq,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
	}
	obj.prepare()
	return obj
}

func (sf *destStationFinder) prepare() {
	req, _ := generateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.oasisReq.Coordinates[1].Lat,
			Lon: sf.oasisReq.Coordinates[1].Lon},
		999,
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
	return iterateNearbyStations(sf.searchResp, sf.searchRespLock)
}
