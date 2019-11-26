package rankingservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
	"github.com/Telenav/osrm-backend/integration/rankingstrategy/rankbyduration"

	"github.com/Telenav/osrm-backend/integration/trafficcache/querytrafficbyedge"
	"github.com/golang/glog"
)

// Handler represents a handler for ranking.
type Handler struct {
	trafficInquirer querytrafficbyedge.Inquirer
	osrmBackend     string
	backendTimeout  time.Duration
}

// New creates a new handler for ranking.
func New(osrmBackend string, backendTimeout time.Duration, trafficInquirer querytrafficbyedge.Inquirer) *Handler {
	if trafficInquirer == nil {
		glog.Fatal("nil traffic inquirer")
		return nil
	}

	return &Handler{
		trafficInquirer,
		osrmBackend,
		backendTimeout,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Handle incoming request %s from remote addr %s", req.RequestURI, req.RemoteAddr)

	// parse incoming request
	osrmRequest, err := osrmv1.ParseRouteRequestURL(req.URL)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v", err)
		return
	}

	// modify
	originalAlternatives := osrmRequest.Alternatives
	originalAnnotations := osrmRequest.Annotations
	osrmRequest.Alternatives = "3" //TODO: re-compile data to support more
	osrmRequest.Annotations = osrmv1.AnnotationsValueTrue

	// route against backend OSRM
	osrmResponse, osrmHTTPStatus, err := h.routeByOSRM(osrmRequest)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(osrmHTTPStatus)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if osrmResponse.Code == osrmv1.CodeOK {
		// update speeds,durations,datasources by traffic
		osrmResponse.Routes = h.updateRoutesByTraffic(osrmResponse.Routes)

		// rank
		osrmResponse.Routes = rankbyduration.Rank(osrmResponse.Routes)

		// pick up
		osrmResponse.Routes = pickupRoutes(osrmResponse.Routes, parseAlternativesNumber(originalAlternatives))

		// cleanup annotations if necessary
		cleanupAnnotations(osrmResponse.Routes, originalAnnotations)
	}

	// return
	w.WriteHeader(osrmHTTPStatus)
	json.NewEncoder(w).Encode(osrmResponse)
}
