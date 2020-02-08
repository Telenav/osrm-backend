package osrmconnector

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
)

// RouteResponse contains osrm Route response and error
type RouteResponse struct {
	Resp *route.Response
	Err  error
}

// TableResponse contains osrm Table response and error
type TableResponse struct {
	Resp *table.Response
	Err  error
}

// OSRMConnector wraps the communication with OSRM server
type OSRMConnector struct {
	osrmClient *osrmHTTPClient
}

// NewOSRMConnector create OsrmConnector object
func NewOSRMConnector(osrmEndpoint string) *OSRMConnector {
	osrm := &OSRMConnector{
		osrmClient: newOsrmHTTPClient(osrmEndpoint),
	}
	go osrm.osrmClient.start()
	return osrm
}

// Request4Route returns a channel immediately.  Response information could be retrieved from the channel when ready.
func (oc *OSRMConnector) Request4Route(r *route.Request) <-chan RouteResponse {
	return oc.osrmClient.submitRouteReq(r)
}

// Request4Table returns a channel immediately.  Response information could be retrieved from the channel when ready.
func (oc *OSRMConnector) Request4Table(r *table.Request) <-chan TableResponse {
	return oc.osrmClient.submitTableReq(r)
}

// Stop will stop OSRMConnector
func (oc *OSRMConnector) Stop() {
	// todo
}
