package oasis

import (
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// GenerateFakeResponse generate response based on oasis request
// It will generate one charge station, use the middle point of given origin/destination's location
func GenerateFakeResponse(req *Request) *Response {
	fakeSolution1 := new(Solution)
	fakeSolution1.Distance = 90000.0
	fakeSolution1.Duration = 30000.0
	fakeSolution1.Weight = 3000.0
	fakeSolution1.RemainingRage = 100000.0
	fakeSolution1.WeightName = "duration"

	fakeStation1 := new(ChargeStation)
	address1 := new(nearbychargestation.Address)
	latMedian := (req.Coordinates[0].Lat + req.Coordinates[1].Lat) / 2
	lonMedian := (req.Coordinates[0].Lon + req.Coordinates[1].Lon) / 2
	address1.GeoCoordinate = nearbychargestation.Coordinate{Latitude: latMedian, Longitude: lonMedian}
	address1.NavCoordinates = append(address1.NavCoordinates, &nearbychargestation.Coordinate{Latitude: latMedian, Longitude: lonMedian})
	fakeStation1.Address = append(fakeStation1.Address, address1)

	fakeStation1.WaitTime = 0.0
	fakeStation1.ChargeTime = 7200.0
	fakeStation1.ChargeRange = req.MaxRange
	fakeStation1.DetailURL = "url"
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation1)

	r := new(Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, fakeSolution1)

	return r
}
