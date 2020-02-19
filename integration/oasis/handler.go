package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
	"github.com/golang/glog"
)

// Handler handles oasis request and provide response
type Handler struct {
	osrmConnector     *osrmconnector.OSRMConnector
	tnSearchConnector *searchconnector.TNSearchConnector
}

// New creates new Handler object
func New(osrmBackend, searchEndpoint, apiKey, apiSignature string) *Handler {
	return &Handler{
		osrmConnector:     osrmconnector.NewOSRMConnector(osrmBackend),
		tnSearchConnector: searchconnector.NewTNSearchConnector(searchEndpoint, apiKey, apiSignature),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse oasis request
	oasisReq, err := oasis.ParseRequestURL(req.URL)
	if err != nil || len(oasisReq.Coordinates) != 2 {
		glog.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	route, err := h.requestRoute4InputOrigDest(oasisReq)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// check whether orig and dest is reachable
	if len(route.Routes) == 0 {
		glog.Info("Orig and destination is not reachable for request.")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Orig and destination is not reachable for request.")
		return
	}

	// check whether has enough energy
	b, remainRange, err := hasEnoughEnergy(oasisReq.CurrRange, oasisReq.SafeLevel, route)
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if b {
		h.generateOASISResponse(w, route, remainRange)
		return
	}

	// check whether could achieve by single charge
	overlap := getOverlapChargeStations4OrigDest(oasisReq, route.Routes[0].Distance, h.osrmConnector, h.tnSearchConnector)
	if len(overlap) > 0 {
		generateResponse4SingleChargeStation(w, oasisReq, overlap, h.osrmConnector)
		return
	} else {
		generateFakeOASISResponse(w, oasisReq)
	}

}

func (h *Handler) requestRoute4InputOrigDest(oasisReq *oasis.Request) (*route.Response, error) {
	// generate route request
	req := route.NewRequest()
	req.Coordinates = oasisReq.Coordinates

	// request for route
	respC := h.osrmConnector.Request4Route(req)

	// retrieve route result
	routeResp := <-respC
	return routeResp.Resp, routeResp.Err
}

func (h *Handler) generateOASISResponse(w http.ResponseWriter, routeResp *route.Response, remainRange float64) {
	w.WriteHeader(http.StatusOK)

	solution := new(oasis.Solution)
	solution.Distance = routeResp.Routes[0].Distance
	solution.Duration = routeResp.Routes[0].Duration
	solution.Weight = routeResp.Routes[0].Weight
	solution.RemainingRage = remainRange
	solution.WeightName = routeResp.Routes[0].WeightName

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, solution)

	json.NewEncoder(w).Encode(r)
}

func generateFakeOASISResponse(w http.ResponseWriter, req *oasis.Request) {
	w.WriteHeader(http.StatusOK)

	fakeSolution1 := new(oasis.Solution)
	fakeSolution1.Distance = 90000.0
	fakeSolution1.Duration = 30000.0
	fakeSolution1.Weight = 3000.0
	fakeSolution1.RemainingRage = 100000.0
	fakeSolution1.WeightName = "duration"

	fakeStation1 := new(oasis.ChargeStation)
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

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, fakeSolution1)

	json.NewEncoder(w).Encode(r)
}
