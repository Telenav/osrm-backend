package osrmconnector

// osrm request type
const (
	OSRMROUTE = iota
	OSRMTABLE = iota
)

type osrmType rune

type request struct {
	url        string
	t          osrmType
	routeRespC chan RouteResponse
	tableRespC chan TableResponse
}

func newOsrmRequest(url string, t osrmType) *request {
	return &request{
		url:        url,
		t:          t,
		routeRespC: make(chan RouteResponse),
		tableRespC: make(chan TableResponse),
	}
}
