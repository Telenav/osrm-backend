package selectionstrategy

import (
	"encoding/json"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/oasis/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/haversine"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationconnquerier"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfinderalg"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationgraph"
	"github.com/golang/glog"

	"github.com/twpayne/go-polyline"
)

// GenerateSolutions4MultipleCharge generates solutions for multiple charge stations needed from orig to dest based on current energy status
// @todo: handle negative situation
func GenerateSolutions4MultipleCharge(w http.ResponseWriter, oasisReq *oasis.Request, routeResp *route.Response, resourceMgr *ResourceMgr) {
	//solutions := generateSolutions4SearchAlongRoute(oasisReq, routeResp, resourceMgr.osrmConnector, resourceMgr.stationFinder)
	solutions := generateSolutions4ChargeStationBasedRoute(oasisReq, resourceMgr)

	w.WriteHeader(http.StatusOK)
	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, solutions...)
	json.NewEncoder(w).Encode(r)
}

func generateSolutions4ChargeStationBasedRoute(oasisReq *oasis.Request, resourceMgr *ResourceMgr) []*oasis.Solution {
	targetSolutions := make([]*oasis.Solution, 0, 10)
	querier := stationconnquerier.New(resourceMgr.spatialIndexerFinder,
		ranker.CreateRanker(ranker.SimpleRanker, resourceMgr.osrmConnector),
		resourceMgr.stationLocationQuerier,
		resourceMgr.connectivityMap,
		&nav.Location{Lat: oasisReq.Coordinates[0].Lat, Lon: oasisReq.Coordinates[0].Lon},
		&nav.Location{Lat: oasisReq.Coordinates[1].Lat, Lon: oasisReq.Coordinates[1].Lon},
		oasisReq.CurrRange,
		oasisReq.MaxRange)
	internalSolutions := stationgraph.NewStationGraph(oasisReq.CurrRange, oasisReq.MaxRange,
		chargingstrategy.NewFakeChargingStrategy(oasisReq.MaxRange),
		querier).GenerateChargeSolutions()

	for _, sol := range internalSolutions {
		targetSolution := sol.Convert2ExternalSolution()
		targetSolutions = append(targetSolutions, targetSolution)
	}

	return targetSolutions
}

func generateSolutions4SearchAlongRoute(oasisReq *oasis.Request, routeResp *route.Response, oc *osrmconnector.OSRMConnector, finder stationfinder.StationFinder) []*oasis.Solution {
	targetSolutions := make([]*oasis.Solution, 0)

	chargeLocations := chargeLocationSelection(oasisReq, routeResp)
	for _, locations := range chargeLocations {
		c := stationfinderalg.CalculateWeightBetweenNeighbors(locations, oc, finder)
		querier := stationfinderalg.NewQuerierBasedOnWeightBetweenNeighborsChan(c)
		internalSolutions := stationgraph.NewStationGraph(oasisReq.CurrRange, oasisReq.MaxRange,
			chargingstrategy.NewFakeChargingStrategy(oasisReq.MaxRange),
			querier).GenerateChargeSolutions()

		for _, sol := range internalSolutions {
			targetSolution := sol.Convert2ExternalSolution()
			targetSolutions = append(targetSolutions, targetSolution)
		}
	}

	return targetSolutions
}

// For each route response, will generate an array of *nav.Location
// Each array contains: start point, first charge point(could also be start point), second charge point, ..., end point
func chargeLocationSelection(oasisReq *oasis.Request, routeResp *route.Response) [][]*nav.Location {
	results := [][]*nav.Location{}
	for _, route := range routeResp.Routes {

		result := []*nav.Location{}
		result = append(result, &nav.Location{
			Lat: oasisReq.Coordinates[0].Lat,
			Lon: oasisReq.Coordinates[0].Lon})
		currEnergy := oasisReq.CurrRange

		// if initial energy is too low
		if currEnergy < oasisReq.PreferLevel {
			result = append(result, &nav.Location{
				Lat: oasisReq.Coordinates[0].Lat,
				Lon: oasisReq.Coordinates[0].Lon})
			currEnergy = oasisReq.MaxRange
		}

		result, currEnergy = findChargeLocation4Route(route, result, currEnergy, oasisReq.PreferLevel, oasisReq.MaxRange)
		if len(result) != 0 {
			result = append(result, &nav.Location{
				Lat: oasisReq.Coordinates[1].Lat,
				Lon: oasisReq.Coordinates[1].Lon})
			results = append(results, result)
		}
	}
	return results
}

func findChargeLocation4Route(route *route.Route, result []*nav.Location, currEnergy, preferLevel, maxRange float64) ([]*nav.Location, float64) {
	for _, leg := range route.Legs {
		for _, step := range leg.Steps {
			if (currEnergy - step.Distance) < preferLevel {
				coords, _, err := polyline.DecodeCoords([]byte(step.Geometry))
				if err != nil {
					glog.Errorf("Incorrect geometry encoding string from route response, error=%v", err)
					return nil, 0.0
				}

				tmp := 0.0
				for i := 0; i < len(coords)-1; i++ {
					tmp += haversine.GreatCircleDistance(coords[i][0], coords[i][1], coords[i+1][0], coords[i+1][1])
					if currEnergy-tmp < preferLevel {
						currEnergy = maxRange
						result = append(result, &nav.Location{
							Lat: coords[i][0],
							Lon: coords[i][1]})

						if currEnergy > (step.Distance - tmp + preferLevel) {
							currEnergy -= step.Distance - tmp
							break
						} else {
							tmp = 0.0
						}
					}
				}

			} else {
				currEnergy -= step.Distance
			}
		}
	}

	return result, currEnergy
}
