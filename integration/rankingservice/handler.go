package rankingservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/code"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route/options"
	"github.com/Telenav/osrm-backend/integration/rankingstrategy/rankbyduration"

	"github.com/Telenav/osrm-backend/integration/trafficcache/querytrafficbyedge"
	"github.com/golang/glog"
)

// Handler represents a handler for ranking.
type Handler struct {
	trafficInquirer querytrafficbyedge.Inquirer
	osrmBackend     string
}

// New creates a new handler for ranking.
func New(osrmBackend string, trafficInquirer querytrafficbyedge.Inquirer) *Handler {
	if trafficInquirer == nil {
		glog.Fatal("nil traffic inquirer")
		return nil
	}

	return &Handler{
		trafficInquirer,
		osrmBackend,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse incoming request
	osrmRequest, err := route.ParseRequestURL(req.URL)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// modify
	originalAlternativesNum := osrmRequest.AlternativesNumber()
	originalAnnotations := osrmRequest.Annotations
	osrmRequest.Alternatives = strconv.FormatUint(uint64(flags.alternatives), 10)
	osrmRequest.Annotations = options.AnnotationsValueTrue

	// route against backend OSRM
	osrmResponse, osrmHTTPStatus, err := h.routeByOSRM(osrmRequest)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(osrmHTTPStatus)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if osrmResponse.Code == code.OK {
		// update speeds,durations,datasources by traffic
		osrmResponse.Routes = h.updateRoutesByTraffic(osrmResponse.Routes)

		// rank
		osrmResponse.Routes = rankbyduration.Rank(osrmResponse.Routes)

		// pick up
		osrmResponse.Routes = pickupRoutes(osrmResponse.Routes, originalAlternativesNum)

		// cleanup annotations if necessary
		cleanupAnnotations(osrmResponse.Routes, originalAnnotations)
	}

	// return
	w.WriteHeader(osrmHTTPStatus)
	json.NewEncoder(w).Encode(osrmResponse)
}
