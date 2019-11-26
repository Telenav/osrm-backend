// Package osrmv1 implements OSRM api v1 in Go code.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md
package osrmv1

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
)

// RouteRequest represent OSRM api v1 route request parameters.
type RouteRequest struct {
	Service     string
	Version     string
	Profile     string
	Coordinates Coordinates

	// Route service query parameters
	Alternatives string
	Steps        bool
	Annotations  string

	//TODO: other parameters

}

// NewRouteRequest create an empty RouteRequest.
func NewRouteRequest() *RouteRequest {
	return &RouteRequest{
		Service:      "route",
		Version:      "v1",
		Profile:      "driving",
		Coordinates:  Coordinates{},
		Alternatives: AlternativesDefaultValue,
		Steps:        StepsDefaultValue,
		Annotations:  AnnotationsDefaultValue,
	}
}

// ParseRouteRequestURI parse Request URI to RouteRequest.
func ParseRouteRequestURI(requestURI string) (*RouteRequest, error) {

	u, err := url.Parse(requestURI)
	if err != nil {
		return nil, err
	}

	return ParseRouteRequestURL(u)
}

// ParseRouteRequestURL parse Request URL to RouteRequest.
func ParseRouteRequestURL(url *url.URL) (*RouteRequest, error) {
	if url == nil {
		return nil, fmt.Errorf("empty URL input")
	}

	routeReq := NewRouteRequest()

	if err := routeReq.parsePath(url.Path); err != nil {
		return nil, err
	}
	routeReq.parseQuery(url.Query())

	return routeReq, nil
}

// QueryValues convert RouteRequest to url.Values.
func (r *RouteRequest) QueryValues() (v url.Values) {

	v = make(url.Values)

	if r.Alternatives != AlternativesDefaultValue {
		v.Add(KeyAlternatives, r.Alternatives)
	}
	if r.Steps != StepsDefaultValue {
		v.Add(KeySteps, strconv.FormatBool(r.Steps))
	}
	if r.Annotations != AnnotationsDefaultValue {
		v.Add(KeyAnnotations, r.Annotations)
	}

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

// AlternativesNumber returns alternatives as number value.
func (r *RouteRequest) AlternativesNumber() int {
	_, n, _ := parseAlternatives(r.Alternatives)
	return n
}

func (r *RouteRequest) pathPrefix() string {
	//i.e. "/route/v1/driving/"
	return api.Slash + r.Service + api.Slash + r.Version + api.Slash + r.Profile + api.Slash
}

func (r *RouteRequest) parsePath(path string) error {
	p := path
	p = strings.TrimPrefix(p, api.Slash)
	p = strings.TrimSuffix(p, api.Slash)

	s := strings.Split(p, api.Slash)
	if len(s) < 4 {
		return fmt.Errorf("invalid path values %v parsed from %s", s, path)
	}
	r.Service = s[0]
	r.Version = s[1]
	r.Profile = s[2]

	var err error
	if r.Coordinates, err = ParseCoordinates(s[3]); err != nil {
		return err
	}

	return nil
}

func (r *RouteRequest) parseQuery(values url.Values) {

	if v := values.Get(KeyAlternatives); len(v) > 0 {
		if alternatives, _, err := parseAlternatives(v); err == nil {
			r.Alternatives = alternatives
		}
	}

	if v := values.Get(KeySteps); len(v) > 0 {
		if b, err := strconv.ParseBool(v); err == nil {
			r.Steps = b
		} else {
			glog.Warning(err)
		}
	}

	if v := values.Get(KeyAnnotations); len(v) > 0 {
		if annotations, err := parseAnnotations(v); err == nil {
			r.Annotations = annotations
		}
	}

}

func parseAlternatives(s string) (string, int, error) {

	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return s, int(n), nil
	}
	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			return s, 2, nil // true : 2
		}
		return s, 1, nil // false : 1
	}

	err := fmt.Errorf("invalid alternatives value: %s", s)
	glog.Warning(err)
	return "", 1, err // use value 1 if fail
}

func parseAnnotations(s string) (string, error) {

	validAnnotationsValues := map[string]struct{}{
		AnnotationsValueTrue:        struct{}{},
		AnnotationsValueFalse:       struct{}{},
		AnnotationsValueNodes:       struct{}{},
		AnnotationsValueDistance:    struct{}{},
		AnnotationsValueDuration:    struct{}{},
		AnnotationsValueDataSources: struct{}{},
		AnnotationsValueWeight:      struct{}{},
		AnnotationsValueSpeed:       struct{}{},
	}

	splits := strings.Split(s, api.Comma)
	for _, split := range splits {
		if _, found := validAnnotationsValues[split]; !found {

			err := fmt.Errorf("invalid annotations value: %s", s)
			glog.Warning(err)
			return "", err
		}
	}

	return s, nil
}
