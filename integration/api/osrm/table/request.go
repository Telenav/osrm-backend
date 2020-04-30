package table

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/genericoptions"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table/options"
	"github.com/golang/glog"
)

// Request for OSRM table service
// http://project-osrm.org/docs/v5.5.1/api/#table-service
type Request struct {
	// Path
	Service     string
	Version     string
	Profile     string
	Coordinates coordinate.Coordinates

	// Options
	Sources      genericoptions.Elements
	Destinations genericoptions.Elements
	Annotations  string
}

// NewRequest create an empty table Request.
func NewRequest() *Request {

	return &Request{
		// Path
		Service:     "table",
		Version:     "v1",
		Profile:     "driving",
		Coordinates: coordinate.Coordinates{},

		// Options
		Sources:      genericoptions.Elements{},
		Destinations: genericoptions.Elements{},
		Annotations:  options.AnnotationsDefaultValue,
	}

}

// ParseRequestURL parse Request URL to Request.
func ParseRequestURL(u *url.URL) (*Request, error) {
	if u == nil {
		return nil, fmt.Errorf("empty URL input")
	}

	req := NewRequest()

	if err := req.parsePath(u.Path); err != nil {
		return nil, err
	}

	req.parseQuery(api.ParseQueryDiscardError(u.RawQuery))

	return req, nil
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

	if v := values.Get(options.KeySources); len(v) > 0 {
		if sources, err := genericoptions.ParseElemenets(v); err == nil {
			r.Sources = sources
		}
	}

	if v := values.Get(options.KeyDestinations); len(v) > 0 {
		if destinations, err := genericoptions.ParseElemenets(v); err == nil {
			r.Destinations = destinations
		}
	}

	if v := values.Get(options.KeyAnnotations); len(v) > 0 {
		if annotations, err := options.ParseAnnotations(v); err == nil {
			r.Annotations = annotations
		}
	}
}

// RequestURI convert TableRequest to RequestURI (e.g. "/path?foo=bar").
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

func (r *Request) pathPrefix() string {
	//i.e. "/table/v1/driving/"
	return api.Slash + r.Service + api.Slash + r.Version + api.Slash + r.Profile + api.Slash
}

// QueryString convert TableRequest to "URL encoded" form ("bar=baz&foo=quux"), but NOT escape.
func (r *Request) QueryString() string {
	rawQuery := r.QueryValues().Encode()
	query, err := url.QueryUnescape(rawQuery)
	if err != nil {
		glog.Warning(err)
		return rawQuery // use rawQuery if unescape fail
	}
	return query
}

// QueryValues convert table Request to url.Values.
func (r *Request) QueryValues() (v url.Values) {
	v = make(url.Values)

	if len(r.Sources) > 0 {
		v.Add(options.KeySources, r.Sources.String())
	}

	if len(r.Destinations) > 0 {
		v.Add(options.KeyDestinations, r.Destinations.String())
	}

	if len(r.Annotations) > 0 {
		v.Add(options.KeyAnnotations, r.Annotations)
	}

	return
}
