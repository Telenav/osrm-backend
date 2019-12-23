// Package route implements OSRM route api v1 in Go code.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md
package route

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route/options"
)

// Request represent OSRM api v1 route request parameters.
type Request struct {
	Service     string
	Version     string
	Profile     string
	Coordinates coordinate.Coordinates

	// Route service query parameters
	Alternatives string
	Steps        bool
	Annotations  string

	//TODO: other parameters

}

// NewRequest create an empty route Request.
func NewRequest() *Request {
	return &Request{
		Service:      "route",
		Version:      "v1",
		Profile:      "driving",
		Coordinates:  coordinate.Coordinates{},
		Alternatives: options.AlternativesDefaultValue,
		Steps:        options.StepsDefaultValue,
		Annotations:  options.AnnotationsDefaultValue,
	}
}

// ParseRequestURI parse Request URI to Request.
func ParseRequestURI(requestURI string) (*Request, error) {

	u, err := url.Parse(requestURI)
	if err != nil {
		return nil, err
	}

	return ParseRequestURL(u)
}

// ParseRequestURL parse Request URL to Request.
func ParseRequestURL(url *url.URL) (*Request, error) {
	if url == nil {
		return nil, fmt.Errorf("empty URL input")
	}

	req := NewRequest()

	if err := req.parsePath(url.Path); err != nil {
		return nil, err
	}
	req.parseQuery(url.Query())

	return req, nil
}

// QueryValues convert route Request to url.Values.
func (r *Request) QueryValues() (v url.Values) {

	v = make(url.Values)

	if r.Alternatives != options.AlternativesDefaultValue {
		v.Add(options.KeyAlternatives, r.Alternatives)
	}
	if r.Steps != options.StepsDefaultValue {
		v.Add(options.KeySteps, strconv.FormatBool(r.Steps))
	}
	if r.Annotations != options.AnnotationsDefaultValue {
		v.Add(options.KeyAnnotations, r.Annotations)
	}

	return
}

// QueryString convert RouteRequest to "URL encoded" form ("bar=baz&foo=quux") .
func (r *Request) QueryString() string {
	return r.QueryValues().Encode()
}

// RequestURI convert RouteRequest to RequestURI (e.g. "/path?foo=bar").
// see more in https://golang.org/pkg/net/url/#URL.RequestURI
func (r *Request) RequestURI() string {
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
func (r *Request) AlternativesNumber() int {
	_, n, _ := options.ParseAlternatives(r.Alternatives)
	return n
}

func (r *Request) pathPrefix() string {
	//i.e. "/route/v1/driving/"
	return api.Slash + r.Service + api.Slash + r.Version + api.Slash + r.Profile + api.Slash
}

func (r *Request) parsePath(path string) error {
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
	if r.Coordinates, err = coordinate.ParseCoordinates(s[3]); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseQuery(values url.Values) {

	if v := values.Get(options.KeyAlternatives); len(v) > 0 {
		if alternatives, _, err := options.ParseAlternatives(v); err == nil {
			r.Alternatives = alternatives
		}
	}

	if v := values.Get(options.KeySteps); len(v) > 0 {
		if b, err := strconv.ParseBool(v); err == nil {
			r.Steps = b
		} else {
			glog.Warning(err)
		}
	}

	if v := values.Get(options.KeyAnnotations); len(v) > 0 {
		if annotations, err := options.ParseAnnotations(v); err == nil {
			r.Annotations = annotations
		}
	}

}
