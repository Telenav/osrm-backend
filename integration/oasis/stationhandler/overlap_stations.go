package stationhandler

import "github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"

type overlapStations struct {
	searchResp *nearbychargestation.Response
}

// func newOverlapStations(iter nearbyStationsIterator, overlap *[]chargeStationInfo) *overlapStations {
// 	searchResp = &nearbychargestation.Response{}
// 	searchResp.Results = make([]*nearbychargestation.Result, len(overlap))
// 	for item := range iter.iterateNearbyStations() {
// 	}
// 	return searchResp
// }
