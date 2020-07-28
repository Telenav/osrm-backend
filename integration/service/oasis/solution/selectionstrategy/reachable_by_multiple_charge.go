package selectionstrategy

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/chargingstrategy"
	"github.com/Telenav/osrm-backend/integration/service/oasis/graph/stationgraph"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/resourcemanager"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationconnquerier"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfinderalg"
	"github.com/Telenav/osrm-backend/integration/util/osrmconnector"
	"github.com/blevesearch/bleve/geo"
	"github.com/golang/glog"

	"github.com/twpayne/go-polyline"
)

// GenerateSolutions4ChargeStationBasedRoute creates optimal charge solution based on charge station graph
func GenerateSolutions4ChargeStationBasedRoute(oasisReq *oasis.Request,
	resourceMgr *resourcemanager.ResourceMgr) ([]*oasis.Solution, error) {

	startTime := time.Now()
	targetSolutions := make([]*oasis.Solution, 0, 10)
	querier := stationconnquerier.New(resourceMgr.SpatialIndexerFinder(),
		ranker.CreateRanker(ranker.SimpleRanker, resourceMgr.OSRMConnector()),
		resourceMgr.StationLocationQuerier(),
		resourceMgr.ConnectivityMap(),
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

	glog.Infof("Generate solutions for charge station based routing takes %f seconds.",
		time.Since(startTime).Seconds())
	return targetSolutions, nil
}

func generateSolutions4SearchAlongRoute(oasisReq *oasis.Request, routeResp *route.Response,
	oc *osrmconnector.OSRMConnector, finder place.IteratorGenerator) []*oasis.Solution {

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
					// geo.Haversin's unit is kilometer, convert to meter
					tmp += geo.Haversin(coords[i][1], coords[i][0], coords[i+1][1], coords[i+1][0]) * 1000
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
