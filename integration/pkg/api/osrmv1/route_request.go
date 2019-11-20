// Package osrmv1 implements OSRM api v1 in Go code.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md
package osrmv1

import (
	"net/url"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
)

// RouteRequest represent OSRM api v1 route request parameters.
type RouteRequest struct {
	Coordinates Coordinates

	//TODO: other parameters
}

// NewRouteRequest create an empty RouteRequest.
func NewRouteRequest() *RouteRequest {
	return &RouteRequest{Coordinates{}}
}

// QueryValues convert RouteRequest to url.Values.
func (r *RouteRequest) QueryValues() (v url.Values) {

	v = make(url.Values)

	// pre-defined parameters in our scenario
	//v.Add("generate_hints", "false")
	//v.Add("overview", "full")
	return
}

// QueryString convert RouteRequest to "URL encoded" form ("bar=baz&foo=quux") .
func (r *RouteRequest) QueryString() string {
	return r.QueryValues().Encode()
}

// RequestURI convert RouteRequest to RequestURI (e.g. "/path?foo=bar").
// see more in https://golang.org/pkg/net/url/#URL.RequestURI
func (r *RouteRequest) RequestURI() string {
	s := r.pathPrefix()

	coordinatesStr := r.Coordinates.String()
	if len(coordinatesStr) > 0 {
		s += coordinatesStr
	}

	queryStr := r.QueryString()
	if len(queryStr) > 0 {
		s += api.QuestionMark + queryStr
	}

	return s
}

func (r *RouteRequest) pathPrefix() string {
	return "/route/v1/driving/"
}
