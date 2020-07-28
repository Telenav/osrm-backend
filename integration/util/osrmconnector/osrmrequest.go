package osrmconnector

// osrm request type
const (
	OSRMRoute = iota
	OSRMTable = iota
)

type requestType int

// request records input information and response channel
type request struct {
	url        string
	t          requestType
	routeRespC chan RouteResponse
	tableRespC chan TableResponse
}

func newOsrmRequest(url string, t requestType) *request {
	return &request{
		url:        url,
		t:          t,
		routeRespC: make(chan RouteResponse),
		tableRespC: make(chan TableResponse),
	}
}
