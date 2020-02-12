package oasis

import (
	"fmt"
	"sort"

	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

// Reachable chargestations from orig already be filterred by currage energy range as radius
// For destination, the filter is a dynamic value, depend on where is the nearest charge station.
// We want to make user has enough energy when reach destination
// The energy level is safeRange + nearest charge station's distance to destination
// If there is one or several charge stations could be found in both origStationsResp and destStationsResp
// We think the result is reachable by single charge station
func isReachableBySingleCharge(req *oasis.Request, origStationsResp, destStationsResp *oasis.ChargeStationsResponse, routeResp *route.Response) bool {
	// only possible when currRange + maxRange > distance + safeRange
	if req.CurrRange+req.MaxRange < routeResp.Routes[0].Distance+req.SafeLevel {
		return false
	}

	if nil == origStationsResp.Resp || len(origStationsResp.Resp) == 0 {
		return false
	}

	if nil == destStationsResp.Resp || len(destStationsResp.Resp) == 0 {
		return false
	}

	// build dict for orig
	origDict := make(map[string]bool)
	for _, station := range origStationsResp.Resp.Results {
		origDict[station.ID] = true
	}

	shared := make([]*nearbychargestation.Result, 10)
	// filter dest stations
	d := destStationsResp.Resp.Results[0].Distance
	d = req.MaxRange - req.SafeLevel - d
	for _, station := range destStationsResp.Resp.Results {
		if _, has := origDict[station.ID]; has {
			shared = append(shared, station)
		}
	}

	// ranking
	if len(shared) != 0 {
		// to do
	}

}

type routePassSingleStation struct {
	time  float64
	index int
}

func rankingSingleChargeStation(orig2Stations, stations2Dest *table.Response, stations []*nearbychargestation.Result) int, error {
	if len(orig2Stations.Durations) != len(stations) || len(stations2Dest.Durations) != len(stations) {
		err := fmt.Errorf("Incorrect parameter for function rankingSingleChargeStation")
		glog.Errorf("%v", err)
		return -1, err
	}

	totalTimes := make([]routePassSingleStation, len(stations))
	for i, _ := range stations {
		var route routePassSingleStation
		route.time = orig2Stations.Durations[i] + stations2Dest.Durations[i]
		route.index = i
		totalTimes = append(totalTimes, route)
	}

	sort.Slice(totalTimes, func(i, j int) bool { return totalTimes[i].time < totalTimes[j].time })
	return totalTimes[0].index, nil
}
