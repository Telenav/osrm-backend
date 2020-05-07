package table

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/Telenav/osrm-backend/integration/api/osrm"
	"github.com/golang/glog"
)

// Request for OSRM table service
// http://project-osrm.org/docs/v5.5.1/api/#table-service
type Request struct {
	// Path
	Service string
	Version string
	Profile string
	osrm.Coordinates

	// Options
	Sources      osrm.OptionElements
	Destinations osrm.OptionElements
	Annotations  string
}

// NewRequest create an empty table Request.
func NewRequest() *Request {

	return &Request{
		// Path
		Service:     "table",
		Version:     "v1",
		Profile:     "driving",
		Coordinates: osrm.Coordinates{},

		// Options
		Sources:      osrm.OptionElements{},
		Destinations: osrm.OptionElements{},
		Annotations:  OptionAnnotationsDefaultValue,
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
	if r.Coordinates, err = osrm.ParseCoordinates(s[3]); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseQuery(values url.Values) {

	if v := values.Get(OptionKeySources); len(v) > 0 {
		if sources, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Sources = sources
		}
	}

	if v := values.Get(OptionKeyDestinations); len(v) > 0 {
		if destinations, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Destinations = destinations
		}
	}

	if v := values.Get(OptionKeyAnnotations); len(v) > 0 {
		if annotations, err := parseOptionAnnotations(v); err == nil {
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
		v.Add(OptionKeySources, r.Sources.String())
	}

	if len(r.Destinations) > 0 {
		v.Add(OptionKeyDestinations, r.Destinations.String())
	}

	if len(r.Annotations) > 0 {
		v.Add(OptionKeyAnnotations, r.Annotations)
	}

	return
}
