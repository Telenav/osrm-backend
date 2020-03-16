package solution

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

// Solution contains summary and selected charge stations
type Solution struct {
	Distance       float64
	Duration       float64
	RemainingRage  float64
	Weight         float64
	ChargeStations []*ChargeStation
}

// ChargeStation contains all information related with specific charge station
type ChargeStation struct {
	Location      Location
	StationID     string
	ArrivalEnergy float64
	WaitTime      float64
	ChargeTime    float64
	ChargeRange   float64
}

// Location defines the geo location of a station
type Location struct {
	Lat float64
	Lon float64
}

// Convert2ExternalSolution convert internal solution format to external format defined in pkg/api/oasis/response
func (sol *Solution) Convert2ExternalSolution() *oasis.Solution {
	target := &oasis.Solution{}

	target.Distance = sol.Distance
	target.Duration = sol.Duration
	target.RemainingRage = sol.RemainingRage

	target.ChargeStations = make([]*oasis.ChargeStation, 0)
	for _, c := range sol.ChargeStations {
		targetStation := &oasis.ChargeStation{}
		targetStation.ChargeRange = c.ChargeRange
		targetStation.ChargeTime = c.ChargeTime
		targetStation.DetailURL = "url"

		targetStation.Address = make([]*nearbychargestation.Address, 0)
		targetAddress := &nearbychargestation.Address{}
		targetAddress.GeoCoordinate = nearbychargestation.Coordinate{
			Latitude:  c.Location.Lat,
			Longitude: c.Location.Lon,
		}
		targetAddress.NavCoordinates = make([]*nearbychargestation.Coordinate, 0)
		targetAddress.NavCoordinates = append(targetAddress.NavCoordinates,
			&nearbychargestation.Coordinate{
				Latitude:  c.Location.Lat,
				Longitude: c.Location.Lon,
			})
		targetStation.Address = append(targetStation.Address, targetAddress)

		target.ChargeStations = append(target.ChargeStations, targetStation)
	}

	return target
}
