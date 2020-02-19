package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
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
func isReachableBySingleCharge(req *oasis.Request, routedistance float64, osrmConnector *osrmconnector.OSRMConnector, tnSearchConnector *searchconnector.TNSearchConnector) coordinate.Coordinates {
	// only possible when currRange + maxRange > distance + safeRange
	if req.CurrRange+req.MaxRange < routedistance+req.SafeLevel {
		return nil
	}

	origStations := stationfinder.NewOrigStationFinder(osrmConnector, tnSearchConnector, req)
	destStations := stationfinder.NewDestStationFinder(osrmConnector, tnSearchConnector, req)
	overlap := stationfinder.FindOverlapBetweenStations(origStations, destStations)

	if len(overlap) == 0 {
		return nil
	}

	overlapPoints := make(coordinate.Coordinates, len(overlap))
	for _, item := range overlap {
		overlapPoints = append(overlapPoints,
			coordinate.Coordinate{
				Lat: item.Location.Lat,
				Lon: item.Location.Lon,
			})
	}
	return overlapPoints
}

type singleChargeStationCandidate struct {
	location         coordinate.Coordinate
	distanceFromOrig float64
	durationFromOrig float64
	distanceToDest   float64
	durationToDest   float64
}

// @todo these logic might refactored later: charge station status calculation should be moved away
func generateResponse4SingleChargeStation(w http.ResponseWriter, req *oasis.Request, overlapPoints coordinate.Coordinates, osrmConnector *osrmconnector.OSRMConnector) {
	candidate, err := pickChargeStationWithEarlistArrival(req, overlapPoints, osrmConnector)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		r := new(oasis.Response)
		r.Message = err.Error()
		json.NewEncoder(w).Encode(r)
		return
	}

	w.WriteHeader(http.StatusOK)

	station := new(oasis.ChargeStation)
	station.WaitTime = 0.0
	station.ChargeTime = 7200.0
	station.ChargeRange = req.MaxRange
	station.DetailURL = "url"
	address := new(nearbychargestation.Address)
	address.GeoCoordinate = nearbychargestation.Coordinate{Latitude: candidate.location.Lat, Longitude: candidate.location.Lon}
	address.NavCoordinates = append(address.NavCoordinates, &nearbychargestation.Coordinate{Latitude: candidate.location.Lat, Longitude: candidate.location.Lon})
	station.Address = append(station.Address, address)

	solution := new(oasis.Solution)
	solution.Distance = candidate.distanceFromOrig + candidate.distanceToDest
	solution.Duration = candidate.durationFromOrig + candidate.durationToDest + station.ChargeTime + station.WaitTime
	solution.RemainingRage = req.MaxRange + req.CurrRange - solution.Distance
	solution.ChargeStations = append(solution.ChargeStations, station)

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, solution)

	json.NewEncoder(w).Encode(r)
}

func pickChargeStationWithEarlistArrival(req *oasis.Request, overlapPoints coordinate.Coordinates, osrmConnector *osrmconnector.OSRMConnector) (*singleChargeStationCandidate, error) {
	if len(overlapPoints) == 0 {
		err := fmt.Errorf("pickChargeStationWithEarlistArrival must be called with none empty overlapPoints")
		glog.Fatalf("%v", err)
		return nil, err
	}

	// request table for orig->overlap stations
	origPoint := coordinate.Coordinates{req.Coordinates[0]}
	reqOrig, _ := osrmhelper.GenerateTableReq4Points(origPoint, overlapPoints)
	respOrigC := osrmConnector.Request4Table(reqOrig)

	// request table for overlap stations -> dest
	destPoint := coordinate.Coordinates{req.Coordinates[1]}
	reqDest, _ := osrmhelper.GenerateTableReq4Points(overlapPoints, destPoint)
	respDestC := osrmConnector.Request4Table(reqDest)

	respOrig := <-respOrigC
	respDest := <-respDestC

	if respOrig.Err != nil {
		glog.Warningf("Table request for url %s failed", reqOrig.RequestURI())
		return nil, respOrig.Err
	}
	if respDest.Err != nil {
		glog.Warningf("Table request for url %s failed", reqDest.RequestURI())
		return nil, respDest.Err
	}

	index, err := rankingSingleChargeStation(respOrig.Resp, respDest.Resp, overlapPoints)
	if err != nil {
		return nil, respDest.Err
	}
	return &singleChargeStationCandidate{
		location:         overlapPoints[index],
		distanceFromOrig: *respOrig.Resp.Distances[0][index],
		durationFromOrig: *respOrig.Resp.Durations[0][index],
		distanceToDest:   *respDest.Resp.Distances[index][0],
		durationToDest:   *respDest.Resp.Durations[index][0],
	}, nil
}

type routePassSingleStation struct {
	time  float64
	index int
}

func rankingSingleChargeStation(orig2Stations, stations2Dest *table.Response, stations coordinate.Coordinates) (int, error) {
	if len(orig2Stations.Durations) != len(stations) || len(stations2Dest.Durations) != len(stations) {
		err := fmt.Errorf("Incorrect parameter for function rankingSingleChargeStation")
		glog.Errorf("%v", err)
		return -1, err
	}

	totalTimes := make([]routePassSingleStation, len(stations))
	for i := range stations {
		var route routePassSingleStation
		route.time = *orig2Stations.Durations[0][i] + *stations2Dest.Durations[i][0]
		route.index = i
		totalTimes = append(totalTimes, route)
	}

	sort.Slice(totalTimes, func(i, j int) bool { return totalTimes[i].time < totalTimes[j].time })

	if totalTimes[0].index < 0 || totalTimes[0].index > len(stations)-1 {
		err := fmt.Errorf("Incorrect index calculated for function rankingSingleChargeStation")
		glog.Errorf("%v", err)
		return totalTimes[0].index, err
	}

	return totalTimes[0].index, nil
}
