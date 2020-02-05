package osrmconnector

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

// RouteResponse contains osrm Route response and error
type RouteResponse struct {
	Resp *route.Response
	Err  error
}

// OsrmConnector wraps the communication with OSRM server
type OsrmConnector struct {
	osrmclient *osrmHTTPClient
}

// NewOsrmConnector create OsrmConnector object
func NewOsrmConnector(osrmEndpoint string) *OsrmConnector {
	osrm := &OsrmConnector{
		osrmclient: newOsrmHTTPClient(osrmEndpoint),
	}
	go osrm.osrmclient.start()
	return osrm
}

// Request4Route provide async api to request for route
func (oc *OsrmConnector) Request4Route(r *route.Request) <-chan RouteResponse {
	return oc.osrmclient.submitRouteReq(r)
}

// Request4Table provide async api to request for table
func (oc *OsrmConnector) Request4Table() <-chan TableResponse {
	// todo
	return nil
}

func (oc *OsrmConnector) Stop() {
	// todo
}

// TableResponse contains osrm Table response and error
type TableResponse struct {
	// todo
}
