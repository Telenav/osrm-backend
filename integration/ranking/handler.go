package ranking

import (
	"fmt"
	"net/http"

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
	//TODO:

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "Not implemented")
}
