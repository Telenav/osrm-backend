package stationfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchhelper"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/searchcoordinate"
	"github.com/golang/glog"
)

const lowEnergyLocationCandidateNumber = 20

type lowEnergyLocationStationFinder struct {
	osrmConnector     *osrmconnector.OSRMConnector
	tnSearchConnector *searchconnector.TNSearchConnector
	location          *StationCoordinate
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
	bf                *basicFinder
}

func NewLowEnergyLocationStationFinder(oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector, location *StationCoordinate) *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		osrmConnector:     oc,
		tnSearchConnector: sc,
		location:          location,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
		bf:                &basicFinder{},
	}
	obj.prepare()
	return obj
}

func (sf *lowEnergyLocationStationFinder) prepare() {
	req, _ := searchhelper.GenerateSearchRequest(
		searchcoordinate.Coordinate{
			Lat: sf.location.Lat,
			Lon: sf.location.Lon},
		lowEnergyLocationCandidateNumber,
		-1)
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

func (sf *lowEnergyLocationStationFinder) iterateNearbyStations() <-chan ChargeStationInfo {
	return sf.bf.iterateNearbyStations(sf.searchResp.Results, sf.searchRespLock)
}
