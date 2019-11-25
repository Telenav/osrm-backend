package ranking

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"

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
	osrmRequest.Alternatives = "3" //TODO: improve
	osrmRequest.Annotations = "true"

	// route against backend OSRM
	osrmResponse, osrmHTTPStatus, err := h.routeByOSRM(osrmRequest)
	if err != nil {
		glog.Warning(err)
		w.WriteHeader(osrmHTTPStatus)
		fmt.Fprintf(w, "%v", err)
		return
	}

	if osrmResponse.Code == osrmv1.CodeOK {
		// ranking and pick up the best one
		osrmResponse.Routes = h.ranking(osrmResponse.Routes, 1)
	}

	// return
	w.WriteHeader(osrmHTTPStatus)
	json.NewEncoder(w).Encode(osrmResponse)
}
