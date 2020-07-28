package clouditerator

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
	"github.com/golang/glog"
)

type basicIterator struct {
	tnSearchConnector *searchconnector.TNSearchConnector
	searchResp        *nearbychargestation.Response
	searchRespLock    *sync.RWMutex
}

func newBasicIterator(sc *searchconnector.TNSearchConnector) *basicIterator {
	return &basicIterator{
		tnSearchConnector: sc,
		searchResp:        nil,
		searchRespLock:    &sync.RWMutex{},
	}
}

func (bf *basicIterator) getNearbyChargeStations(req *nearbychargestation.Request) {
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
func (bf *basicIterator) IterateNearbyStations() <-chan *iteratortype.ChargeStationInfo {
	if bf.searchResp == nil || len(bf.searchResp.Results) == 0 {
		c := make(chan *iteratortype.ChargeStationInfo)
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

	c := make(chan *iteratortype.ChargeStationInfo, size)
	go func() {
		defer close(c)
		for _, result := range results {
			if len(result.Place.Address) == 0 {
				continue
			}
			station := iteratortype.ChargeStationInfo{
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
