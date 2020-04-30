package cloudfinder

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
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

// NearbyStationsIterator provides channel which contains near by station information
func (bf *basicFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	if bf.searchResp == nil || len(bf.searchResp.Results) == 0 {
		c := make(chan *stationfindertype.ChargeStationInfo)
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

	c := make(chan *stationfindertype.ChargeStationInfo, size)
	go func() {
		defer close(c)
		for _, result := range results {
			if len(result.Place.Address) == 0 {
				continue
			}
			station := stationfindertype.ChargeStationInfo{
				ID: result.ID,
				Location: nav.Location{
					Lat: result.Place.Address[0].GeoCoordinate.Latitude,
					Lon: result.Place.Address[0].GeoCoordinate.Longitude},
			}
			c <- &station
		}
	}()

	return c
}
