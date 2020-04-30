// Package route implements OSRM route api v1 in Go code.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md
package route

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api/osrm/genericoptions"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/Telenav/osrm-backend/integration/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/api/osrm/route/options"
)

// Request represent OSRM api v1 route request parameters.
type Request struct {

	// Path
	Service     string
	Version     string
	Profile     string
	Coordinates coordinate.Coordinates

	// generic parameters
	Bearings      genericoptions.Elements
	Radiuses      genericoptions.Elements
	GenerateHints bool
	Hints         genericoptions.Elements
	Approaches    genericoptions.Elements
	Exclude       genericoptions.Classes

	// Route service query parameters
	Alternatives     string
	Steps            bool
	Annotations      string
	Geometries       string
	Overview         string
	ContinueStraight string
	Waypoints        coordinate.Indexes
}

// NewRequest create an empty route Request.
func NewRequest() *Request {
	return &Request{
		// Path
		Service:     "route",
		Version:     "v1",
		Profile:     "driving",
		Coordinates: coordinate.Coordinates{},

		// generic options
		Bearings:      genericoptions.Elements{},
		Radiuses:      genericoptions.Elements{},
		GenerateHints: genericoptions.GenerateHintsDefaultValue,
		Hints:         genericoptions.Elements{},
		Approaches:    genericoptions.Elements{},
		Exclude:       genericoptions.Classes{},

		// route options
		Alternatives:     options.AlternativesDefaultValue,
		Steps:            options.StepsDefaultValue,
		Annotations:      options.AnnotationsDefaultValue,
		Geometries:       options.GeometriesDefaultValue,
		Overview:         options.OverviewDefaultValue,
		ContinueStraight: options.ContinueStraightDefaultValue,
		Waypoints:        coordinate.Indexes{},
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
func ParseRequestURL(u *url.URL) (*Request, error) {
	if u == nil {
		return nil, fmt.Errorf("empty URL input")
	}

	req := NewRequest()

	if err := req.parsePath(u.Path); err != nil {
		return nil, err
	}

	//NOTE: url.Query() will also use ";" as seprator, which is not expected. So we implements our own version instead.
	//req.parseQuery(u.Query())
	req.parseQuery(api.ParseQueryDiscardError(u.RawQuery))

	return req, nil
}

// QueryValues convert route Request to url.Values.
func (r *Request) QueryValues() (v url.Values) {

	v = make(url.Values)

	// generic options
	if len(r.Bearings) > 0 {
		v.Add(genericoptions.KeyBearings, r.Bearings.String())
	}
	if len(r.Radiuses) > 0 {
		v.Add(genericoptions.KeyRadiuses, r.Radiuses.String())
	}
	if r.GenerateHints != genericoptions.GenerateHintsDefaultValue {
		v.Add(genericoptions.KeyGenerateHints, strconv.FormatBool(r.GenerateHints))
	}
	if len(r.Hints) > 0 {
		v.Add(genericoptions.KeyHints, r.Hints.String())
	}
	if len(r.Approaches) > 0 {
		v.Add(genericoptions.KeyApproaches, r.Approaches.String())
	}
	if len(r.Exclude) > 0 {
		v.Add(genericoptions.KeyExclude, r.Exclude.String())
	}

	// route options
	if r.Alternatives != options.AlternativesDefaultValue {
		v.Add(options.KeyAlternatives, r.Alternatives)
	}
	if r.Steps != options.StepsDefaultValue {
		v.Add(options.KeySteps, strconv.FormatBool(r.Steps))
	}
	if r.Annotations != options.AnnotationsDefaultValue {
		v.Add(options.KeyAnnotations, r.Annotations)
	}
	if r.Geometries != options.GeometriesDefaultValue {
		v.Add(options.KeyGeometries, r.Geometries)
	}
	if r.Overview != options.OverviewDefaultValue {
		v.Add(options.KeyOverview, r.Overview)
	}
	if r.ContinueStraight != options.ContinueStraightDefaultValue {
		v.Add(options.KeyContinueStraight, r.ContinueStraight)
	}
	if len(r.Waypoints) > 0 {
		v.Add(options.KeyWaypoints, r.Waypoints.String())
	}

	return
}

// QueryString convert RouteRequest to "URL encoded" form ("bar=baz&foo=quux"), but NOT escape.
func (r *Request) QueryString() string {
	rawQuery := r.QueryValues().Encode()
	query, err := url.QueryUnescape(rawQuery)
	if err != nil {
		glog.Warning(err)
		return rawQuery // use rawQuery if unescape fail
	}
	return query
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

	if v := values.Get(genericoptions.KeyBearings); len(v) > 0 {
		if bearings, err := genericoptions.ParseElemenets(v); err == nil {
			r.Bearings = bearings
		}
	}
	if v := values.Get(genericoptions.KeyRadiuses); len(v) > 0 {
		if radiuses, err := genericoptions.ParseElemenets(v); err == nil {
			r.Radiuses = radiuses
		}
	}
	if v := values.Get(genericoptions.KeyGenerateHints); len(v) > 0 {
		if generateHints, err := genericoptions.ParseGenerateHints(v); err == nil {
			r.GenerateHints = generateHints
		}
	}
	if v := values.Get(genericoptions.KeyHints); len(v) > 0 {
		if hints, err := genericoptions.ParseElemenets(v); err == nil {
			r.Hints = hints
		}
	}
	if v := values.Get(genericoptions.KeyApproaches); len(v) > 0 {
		if approaches, err := genericoptions.ParseElemenets(v); err == nil {
			r.Approaches = approaches
		}
	}
	if v := values.Get(genericoptions.KeyExclude); len(v) > 0 {
		if classes, err := genericoptions.ParseClasses(v); err == nil {
			r.Exclude = classes
		}
	}

	if v := values.Get(options.KeyAlternatives); len(v) > 0 {
		if alternatives, _, err := options.ParseAlternatives(v); err == nil {
			r.Alternatives = alternatives
		}
	}
	if v := values.Get(options.KeySteps); len(v) > 0 {
		if b, err := options.ParseSteps(v); err == nil {
			r.Steps = b
		}
	}
	if v := values.Get(options.KeyAnnotations); len(v) > 0 {
		if annotations, err := options.ParseAnnotations(v); err == nil {
			r.Annotations = annotations
		}
	}
	if v := values.Get(options.KeyGeometries); len(v) > 0 {
		if geometries, err := options.ParseGeometries(v); err == nil {
			r.Geometries = geometries
		}
	}
	if v := values.Get(options.KeyOverview); len(v) > 0 {
		if overview, err := options.ParseOverview(v); err == nil {
			r.Overview = overview
		}
	}
	if v := values.Get(options.KeyContinueStraight); len(v) > 0 {
		if continueStraight, err := options.ParseContinueStraight(v); err == nil {
			r.ContinueStraight = continueStraight
		}
	}
	if v := values.Get(options.KeyWaypoints); len(v) > 0 {
		if indexes, err := coordinate.PraseIndexes(v); err == nil {
			r.Waypoints = indexes
		}
	}

}
