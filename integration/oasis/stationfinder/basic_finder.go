package stationfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/nav"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

type basicFinder struct {
	tnSearchConnector *searchconnector.TNSearchConnector
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
}

func newBasicFinder(sc *searchconnector.TNSearchConnector) *basicFinder {
	return &basicFinder{
		tnSearchConnector: sc,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
	}
}

func (bf *basicFinder) getNearbyChargeStations(req *nearbychargestation.Request) {
	respC := bf.tnSearchConnector.ChargeStationSearch(req)
	resp := <-respC
	if resp.Err != nil {
		glog.Warningf("Search failed during prepare orig search for url: %s", req.RequestURI())
		return
	}

	bf.searchRespLock.Lock()
	bf.searchResp = resp.Resp
	bf.searchRespLock.Unlock()
}

func (bf *basicFinder) iterateNearbyStations() <-chan ChargeStationInfo {
	if bf.searchResp == nil || len(bf.searchResp.Results) == 0 {
		c := make(chan ChargeStationInfo)
		go func() {
			defer close(c)
		}()
		return c
	}

	bf.searchRespLock.RLock()
	size := len(bf.searchResp.Results)
	results := make([]*nearbychargestation.Result, size)
	copy(results, bf.searchResp.Results)
	bf.searchRespLock.RUnlock()

	c := make(chan ChargeStationInfo, size)
	go func() {
		defer close(c)
		for _, result := range results {
			if len(result.Place.Address) == 0 {
				continue
			}
			station := ChargeStationInfo{
				ID: result.ID,
				Location: nav.Location{
					Lat: result.Place.Address[0].GeoCoordinate.Latitude,
					Lon: result.Place.Address[0].GeoCoordinate.Longitude},
			}
			c <- station
		}
	}()

	return c
}
