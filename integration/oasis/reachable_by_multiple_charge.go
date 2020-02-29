package oasis

import (
	"github.com/Telenav/osrm-backend/integration/oasis/haversine"
	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/stationfinder"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/golang/glog"

	"github.com/twpayne/go-polyline"
)

func pickChargeStationsWithEarlistArrival(oasisReq *oasis.Request, routeResp *route.Response, oc *osrmconnector.OSRMConnector, sc *searchconnector.TNSearchConnector) {
	chargeLocations := chargeLocationSelection(oasisReq, routeResp)

	for _, locations := range chargeLocations {

	}
}

// For each route response, will generate an array of *stationfinder.StationCoordinate
// Each array contains: start point, first charge point(could also be start point), second charge point, ..., end point
func chargeLocationSelection(oasisReq *oasis.Request, routeResp *route.Response) [][]*stationfinder.StationCoordinate {
	results := [][]*stationfinder.StationCoordinate{}
	for _, route := range routeResp.Routes {

		result := []*stationfinder.StationCoordinate{}
		result = append(result, &stationfinder.StationCoordinate{
			Lat: oasisReq.Coordinates[0].Lat,
			Lon: oasisReq.Coordinates[0].Lon})
		currEnergy := oasisReq.CurrRange

		// if initial energy is too low
		if currEnergy < oasisReq.PreferLevel {
			result = append(result, &stationfinder.StationCoordinate{
				Lat: oasisReq.Coordinates[0].Lat,
				Lon: oasisReq.Coordinates[0].Lon})
			currEnergy = oasisReq.MaxRange
		}

		result, currEnergy = findChargeLocation4Route(route, result, currEnergy, oasisReq.PreferLevel, oasisReq.MaxRange)
		if len(result) != 0 {
			result = append(result, &stationfinder.StationCoordinate{
				Lat: oasisReq.Coordinates[1].Lat,
				Lon: oasisReq.Coordinates[1].Lon})
			results = append(results, result)
		}
	}
	return results
}

func findChargeLocation4Route(route *route.Route, result []*stationfinder.StationCoordinate, currEnergy, preferLevel, maxRange float64) ([]*stationfinder.StationCoordinate, float64) {
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
						result = append(result, &stationfinder.StationCoordinate{
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
