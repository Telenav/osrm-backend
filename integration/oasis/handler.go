package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

type Handler struct {
	osrmBackend string
}

func New(osrmBackend string) *Handler {
	return &Handler{
		osrmBackend,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)
	w.WriteHeader(http.StatusOK)

	oasisRequest, err := oasis.ParseRequestURL(req.URL)
	if err != nil || len(oasisRequest.Coordinates) != 2 {
		glog.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	json.NewEncoder(w).Encode(generateFakeOasisResponseBasedOnRequest(oasisRequest))
}

func generateFakeOasisResponse() *oasis.Response {
	// Fake solution 1
	fakeSolution1 := new(oasis.Solution)
	fakeSolution1.Distance = 90.0
	fakeSolution1.Duration = 300.0
	fakeSolution1.Weight = 300.0
	fakeSolution1.RemainingRage = 100000.0
	fakeSolution1.WeightName = "duration"

	// Information realted with first charge station
	fakeStation1 := new(oasis.ChargeStation)
	address1 := new(nearbychargestation.Address)
	address1.GeoCoordinate = nearbychargestation.Coordinate{Latitude: 37.78509, Longitude: -122.41988}
	address1.NavCoordinates = append(address1.NavCoordinates, &nearbychargestation.Coordinate{Latitude: 37.78509, Longitude: -122.41988})
	fakeStation1.Address = append(fakeStation1.Address, address1)
	fakeStation1.WaitTime = 30.0
	fakeStation1.ChargeTime = 100.0
	fakeStation1.ChargeRange = 100.0
	fakeStation1.DetailURL = "url"
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation1)

	// Information realted with second charge station
	fakeStation2 := new(oasis.ChargeStation)
	address2 := new(nearbychargestation.Address)
	address2.GeoCoordinate = nearbychargestation.Coordinate{Latitude: 13.40677, Longitude: 52.53333}
	address2.NavCoordinates = append(address2.NavCoordinates, &nearbychargestation.Coordinate{Latitude: 13.40677, Longitude: 52.53333})
	fakeStation2.Address = append(fakeStation2.Address, address2)
	fakeStation2.WaitTime = 100.0
	fakeStation2.ChargeTime = 100.0
	fakeStation2.ChargeRange = 100.0
	fakeStation2.DetailURL = "url"
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation2)

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success"
	r.Solutions = append(r.Solutions, fakeSolution1)

	return r
}

func generateFakeOasisResponseBasedOnRequest(req *oasis.Request) *oasis.Response {
	fakeSolution1 := new(oasis.Solution)
	fakeSolution1.Distance = 90000.0
	fakeSolution1.Duration = 30000.0
	fakeSolution1.Weight = 3000.0
	fakeSolution1.RemainingRage = 100000.0
	fakeSolution1.WeightName = "duration"

	fakeStation1 := new(oasis.ChargeStation)
	address1 := new(nearbychargestation.Address)
	LatMedian := (req.Coordinates[0].Lat + req.Coordinates[1].Lat) / 2
	LonMedian := (req.Coordinates[0].Lon + req.Coordinates[1].Lon) / 2
	address1.GeoCoordinate = nearbychargestation.Coordinate{Latitude: LatMedian, Longitude: LonMedian}
	address1.NavCoordinates = append(address1.NavCoordinates, &nearbychargestation.Coordinate{Latitude: LatMedian, Longitude: LonMedian})
	fakeStation1.Address = append(fakeStation1.Address, address1)

	fakeStation1.WaitTime = 0.0
	fakeStation1.ChargeTime = 7200.0
	fakeStation1.ChargeRange = req.MaxRange
	fakeStation1.DetailURL = "url"
	fakeSolution1.ChargeStations = append(fakeSolution1.ChargeStations, fakeStation1)

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, fakeSolution1)
	return r
}
