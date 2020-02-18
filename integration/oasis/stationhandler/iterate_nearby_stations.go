package stationhandler

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

type basicFinder struct {
}

func iterateNearbyStations(searchResp *nearbychargestation.Response, respLock *sync.RWMutex) <-chan chargeStationInfo {
	if searchResp == nil || len(searchResp.Results) == 0 {
		c := make(chan chargeStationInfo)
		go func() {
			defer close(c)
		}()
		return c
	}

	c := make(chan chargeStationInfo, len(searchResp.Results))
	results := make([]*nearbychargestation.Result, len(searchResp.Results))

	if respLock != nil {
		respLock.RLock()
	}
	copy(results, searchResp.Results)
	if respLock != nil {
		respLock.RUnlock()
	}

	go func() {
		defer close(c)
		for _, result := range results {
			if len(result.Place.Address) == 0 {
				continue
			}
			station := chargeStationInfo{
				id: result.ID,
				location: coordinate.Coordinate{
					Lat: result.Place.Address[0].GeoCoordinate.Latitude,
					Lon: result.Place.Address[0].GeoCoordinate.Longitude},
			}
			c <- station
		}
	}()

	return c
}
