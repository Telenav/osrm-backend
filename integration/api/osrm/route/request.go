// Package route implements OSRM route api v1 in Go code.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md
package route

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api/osrm"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/api"
)

// Request represent OSRM api v1 route request parameters.
type Request struct {

	// Path
	Service string
	Version string
	Profile string
	osrm.Coordinates

	// generic parameters
	Bearings      osrm.OptionElements
	Radiuses      osrm.OptionElements
	GenerateHints bool
	Hints         osrm.OptionElements
	Approaches    osrm.OptionElements
	Exclude       osrm.OptionClasses

	// Route service query parameters
	Alternatives     string
	Steps            bool
	Annotations      string
	Geometries       string
	Overview         string
	ContinueStraight string
	Waypoints        osrm.CoordinateIndexes
}

// NewRequest create an empty route Request.
func NewRequest() *Request {
	return &Request{
		// Path
		Service:     "route",
		Version:     "v1",
		Profile:     "driving",
		Coordinates: osrm.Coordinates{},

		// generic options
		Bearings:      osrm.OptionElements{},
		Radiuses:      osrm.OptionElements{},
		GenerateHints: osrm.OptionGenerateHintsDefaultValue,
		Hints:         osrm.OptionElements{},
		Approaches:    osrm.OptionElements{},
		Exclude:       osrm.OptionClasses{},

		// route options
		Alternatives:     OptionAlternativesDefaultValue,
		Steps:            OptionStepsDefaultValue,
		Annotations:      OptionAnnotationsDefaultValue,
		Geometries:       OptionGeometriesDefaultValue,
		Overview:         OptionOverviewDefaultValue,
		ContinueStraight: OptionContinueStraightDefaultValue,
		Waypoints:        osrm.CoordinateIndexes{},
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
		v.Add(osrm.OptionKeyBearings, r.Bearings.String())
	}
	if len(r.Radiuses) > 0 {
		v.Add(osrm.OptionKeyRadiuses, r.Radiuses.String())
	}
	if r.GenerateHints != osrm.OptionGenerateHintsDefaultValue {
		v.Add(osrm.OptionKeyGenerateHints, strconv.FormatBool(r.GenerateHints))
	}
	if len(r.Hints) > 0 {
		v.Add(osrm.OptionKeyHints, r.Hints.String())
	}
	if len(r.Approaches) > 0 {
		v.Add(osrm.OptionKeyApproaches, r.Approaches.String())
	}
	if len(r.Exclude) > 0 {
		v.Add(osrm.OptionKeyExclude, r.Exclude.String())
	}

	// route options
	if r.Alternatives != OptionAlternativesDefaultValue {
		v.Add(OptionKeyAlternatives, r.Alternatives)
	}
	if r.Steps != OptionStepsDefaultValue {
		v.Add(OptionKeySteps, strconv.FormatBool(r.Steps))
	}
	if r.Annotations != OptionAnnotationsDefaultValue {
		v.Add(OptionKeyAnnotations, r.Annotations)
	}
	if r.Geometries != OptionGeometriesDefaultValue {
		v.Add(OptionKeyGeometries, r.Geometries)
	}
	if r.Overview != OptionOverviewDefaultValue {
		v.Add(OptionKeyOverview, r.Overview)
	}
	if r.ContinueStraight != OptionContinueStraightDefaultValue {
		v.Add(OptionKeyContinueStraight, r.ContinueStraight)
	}
	if len(r.Waypoints) > 0 {
		v.Add(OptionKeyWaypoints, r.Waypoints.String())
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
	_, n, _ := parseOptionAlternatives(r.Alternatives)
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
	if r.Coordinates, err = osrm.ParseCoordinates(s[3]); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseQuery(values url.Values) {

	if v := values.Get(osrm.OptionKeyBearings); len(v) > 0 {
		if bearings, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Bearings = bearings
		}
	}
	if v := values.Get(osrm.OptionKeyRadiuses); len(v) > 0 {
		if radiuses, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Radiuses = radiuses
		}
	}
	if v := values.Get(osrm.OptionKeyGenerateHints); len(v) > 0 {
		if generateHints, err := osrm.ParseOptionGenerateHints(v); err == nil {
			r.GenerateHints = generateHints
		}
	}
	if v := values.Get(osrm.OptionKeyHints); len(v) > 0 {
		if hints, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Hints = hints
		}
	}
	if v := values.Get(osrm.OptionKeyApproaches); len(v) > 0 {
		if approaches, err := osrm.ParseOptionElemenets(v); err == nil {
			r.Approaches = approaches
		}
	}
	if v := values.Get(osrm.OptionKeyExclude); len(v) > 0 {
		if classes, err := osrm.ParseOptionClasses(v); err == nil {
			r.Exclude = classes
		}
	}

	if v := values.Get(OptionKeyAlternatives); len(v) > 0 {
		if alternatives, _, err := parseOptionAlternatives(v); err == nil {
			r.Alternatives = alternatives
		}
	}
	if v := values.Get(OptionKeySteps); len(v) > 0 {
		if b, err := parseOptionSteps(v); err == nil {
			r.Steps = b
		}
	}
	if v := values.Get(OptionKeyAnnotations); len(v) > 0 {
		if annotations, err := parseOptionAnnotations(v); err == nil {
			r.Annotations = annotations
		}
	}
	if v := values.Get(OptionKeyGeometries); len(v) > 0 {
		if geometries, err := parseOptionGeometries(v); err == nil {
			r.Geometries = geometries
		}
	}
	if v := values.Get(OptionKeyOverview); len(v) > 0 {
		if overview, err := parseOptionOverview(v); err == nil {
			r.Overview = overview
		}
	}
	if v := values.Get(OptionKeyContinueStraight); len(v) > 0 {
		if continueStraight, err := parseOptionContinueStraight(v); err == nil {
			r.ContinueStraight = continueStraight
		}
	}
	if v := values.Get(OptionKeyWaypoints); len(v) > 0 {
		if indexes, err := osrm.PraseCoordinateIndexes(v); err == nil {
			r.Waypoints = indexes
		}
	}

}
