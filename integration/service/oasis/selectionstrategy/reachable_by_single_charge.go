package selectionstrategy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/osrm"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/osrmhelper"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfinderalg"
	"github.com/golang/glog"
)

const maxOverlapPointsNum = 500

// Reachable chargestations from orig already be filterred by currage energy range as radius
// For destination, the filter is a dynamic value, depend on where is the nearest charge station.
// We want to make user has enough energy when reach destination
// The energy level is safeRange + nearest charge station's distance to destination
// If there is one or several charge stations could be found in both origStationsResp and destStationsResp
// We think the result is reachable by single charge station
func GetOverlapChargeStations4OrigDest(req *oasis.Request, routedistance float64, resourceMgr *ResourceMgr) osrm.Coordinates {
	// only possible when currRange + maxRange > distance + safeRange
	if req.CurrRange+req.MaxRange < routedistance+req.SafeLevel {
		return nil
	}

	origStations := resourceMgr.stationFinder.NewOrigStationFinder(req)
	destStations := resourceMgr.stationFinder.NewDestStationFinder(req)
	overlap := stationfinderalg.FindOverlapBetweenStations(origStations, destStations)

	if len(overlap) == 0 {
		return nil
	}

	var overlapPoints osrm.Coordinates
	for i, item := range overlap {
		overlapPoints = append(overlapPoints,
			osrm.Coordinate{
				Lat: item.Location.Lat,
				Lon: item.Location.Lon,
			})
		if i > maxOverlapPointsNum {
			break
		}
	}
	return overlapPoints
}

type singleChargeStationCandidate struct {
	location         osrm.Coordinate
	distanceFromOrig float64
	durationFromOrig float64
	distanceToDest   float64
	durationToDest   float64
}

// GenerateResponse4SingleChargeStation generates response for only single charge station needed from orig to dest based on current energy status
func GenerateResponse4SingleChargeStation(w http.ResponseWriter, req *oasis.Request, overlapPoints osrm.Coordinates, resourceMgr *ResourceMgr) bool {
	candidate, err := pickChargeStationWithEarlistArrival(req, overlapPoints, resourceMgr.osrmConnector)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		r := new(oasis.Response)
		r.Message = err.Error()
		json.NewEncoder(w).Encode(r)
		return false
	}

	w.WriteHeader(http.StatusOK)

	station := new(oasis.ChargeStation)
	station.WaitTime = 0.0
	// @todo ChargeTime and ChargeRange need to be adjusted according to chargingstrategy
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
	return true
}

func pickChargeStationWithEarlistArrival(req *oasis.Request, overlapPoints osrm.Coordinates, osrmConnector *osrmconnector.OSRMConnector) (*singleChargeStationCandidate, error) {
	if len(overlapPoints) == 0 {
		err := fmt.Errorf("pickChargeStationWithEarlistArrival must be called with none empty overlapPoints")
		glog.Fatalf("%v", err)
		return nil, err
	}

	// request table for orig->overlap stations
	origPoint := osrm.Coordinates{req.Coordinates[0]}
	reqOrig, _ := osrmhelper.GenerateTableReq4Points(origPoint, overlapPoints)
	respOrigC := osrmConnector.Request4Table(reqOrig)

	// request table for overlap stations -> dest
	destPoint := osrm.Coordinates{req.Coordinates[1]}
	reqDest, _ := osrmhelper.GenerateTableReq4Points(overlapPoints, destPoint)
	respDestC := osrmConnector.Request4Table(reqDest)

	respOrig := <-respOrigC
	respDest := <-respDestC

	if respOrig.Err != nil {
		glog.Warningf("Table request failed for url %s with error %v", reqOrig.RequestURI(), respOrig.Err)
		return nil, respOrig.Err
	}
	if respDest.Err != nil {
		glog.Warningf("Table request failed for url %s with error %v", reqDest.RequestURI(), respDest.Err)
		return nil, respDest.Err
	}
	if len(respOrig.Resp.Durations[0]) != len(respDest.Resp.Durations) || len(overlapPoints) != len(respOrig.Resp.Durations[0]) {
		err := fmt.Errorf("Incorrect table response, the dimension of array is not as expected. [orig2overlap, overlap2dest, overlap]= %d, %d, %d",
			len(respOrig.Resp.Durations[0]), len(respDest.Resp.Durations), len(overlapPoints))
		glog.Errorf("%v", err)
		return nil, err
	}

	index, err := rankingSingleChargeStation(respOrig.Resp, respDest.Resp)
	if err != nil {
		return nil, err
	}
	return &singleChargeStationCandidate{
		location:         overlapPoints[index],
		distanceFromOrig: respOrig.Resp.Distances[0][index],
		durationFromOrig: respOrig.Resp.Durations[0][index],
		distanceToDest:   respDest.Resp.Distances[index][0],
		durationToDest:   respDest.Resp.Durations[index][0],
	}, nil
}

type routePassSingleStation struct {
	time  float64
	index int
}

func rankingSingleChargeStation(orig2Stations, stations2Dest *table.Response) (int, error) {
	if len(orig2Stations.Durations) == 0 || len(orig2Stations.Durations[0]) != len(stations2Dest.Durations) {
		err := fmt.Errorf("Incorrect table response for function rankingSingleChargeStation")
		glog.Errorf("%v", err)
		return -1, err
	}

	size := len(orig2Stations.Durations[0])

	var totalTimes []routePassSingleStation
	for i := 0; i < size; i++ {
		var route routePassSingleStation
		route.time = orig2Stations.Durations[0][i] + stations2Dest.Durations[i][0]
		route.index = i
		totalTimes = append(totalTimes, route)
	}

	sort.Slice(totalTimes, func(i, j int) bool { return totalTimes[i].time < totalTimes[j].time })

	return totalTimes[0].index, nil
}
