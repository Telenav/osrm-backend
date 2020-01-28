package oasis

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis/options"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/golang/glog"
)

// Request for oasis service
type Request struct {
	Service     string
	Version     string
	Profile     string
	Coordinates coordinate.Coordinates

	MaxRange    float64
	CurrRange   float64
	PreferLevel float64
	SafeLevel   float64
}

// NewRequest create an empty oasis Request.
func NewRequest() *Request {
	return &Request{
		// Path
		Service:     "oasis",
		Version:     "v1",
		Profile:     "earliest",
		Coordinates: coordinate.Coordinates{},

		// generic options
		MaxRange:    options.InvalidMaxRangeValue,
		CurrRange:   options.InvalidCurrentRangeValue,
		PreferLevel: options.DefaultPreferLevel,
		SafeLevel:   options.DefaultSafeLevel,
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

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}

	if err := req.parseQuery(params); err != nil {
		glog.Warning(err)
		return nil, err
	}

	if err := req.validate(); err != nil {
		glog.Warning(err)
		return nil, err
	}

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

func (r *Request) parseQuery(params url.Values) error {
	for k, v := range params {
		if len(v) <= 0 {
			continue
		}

		var valfloat float64
		if s, err := strconv.ParseFloat(v[0], 64); err != nil {
			return err
		} else {
			valfloat = float64(s)
		}

		switch k {
		case options.KeyMaxRange:
			r.MaxRange = valfloat
		case options.KeyCurrRange:
			r.CurrRange = valfloat
		case options.KeyPreferLevel:
			r.PreferLevel = valfloat
		case options.KeySafeLevel:
			r.SafeLevel = valfloat
		}
	}

	return nil
}

func (r *Request) validate() error {
	// MaxRange must be set
	if floatEquals(r.MaxRange, options.InvalidMaxRangeValue) || r.MaxRange < 0 {
		return errors.New("Invalid value for " + options.KeyMaxRange + ".")
	}

	// CurrRange must be set
	if floatEquals(r.CurrRange, options.InvalidCurrentRangeValue) || r.CurrRange < 0 {
		return errors.New("Invalid value for " + options.KeyCurrRange + ".")
	}

	// CurrRange must be smaller or equal to MaxRange
	if r.CurrRange > r.MaxRange {
		return errors.New(options.KeyCurrRange + " must be smaller or equal to " + options.KeyMaxRange + ".")
	}

	// PreferLevel must be smaller to MaxRange
	if r.PreferLevel >= r.MaxRange {
		return errors.New(options.KeyPreferLevel + " must be smaller to " + options.KeyMaxRange + ".")
	}

	// SafeLevel must be smaller to MaxRange
	if r.SafeLevel >= r.MaxRange {
		return errors.New(options.KeySafeLevel + " must be smaller to " + options.KeyMaxRange + ".")
	}

	return nil
}

var epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < epsilon && (b-a) < epsilon {
		return true
	}
	return false
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

func (r *Request) pathPrefix() string {
	//i.e. "/oasis/v1/earliest/"
	return api.Slash + r.Service + api.Slash + r.Version + api.Slash + r.Profile + api.Slash
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

// QueryValues convert route Request to url.Values.
func (r *Request) QueryValues() (v url.Values) {
	v = make(url.Values)

	if r.MaxRange != 0 {
		v.Add(options.KeyMaxRange, strconv.FormatFloat(r.MaxRange, 'f', -1, 64))
	}

	if r.CurrRange != 0 {
		v.Add(options.KeyCurrRange, strconv.FormatFloat(r.CurrRange, 'f', -1, 64))
	}

	if r.PreferLevel != 0 {
		v.Add(options.KeyPreferLevel, strconv.FormatFloat(r.PreferLevel, 'f', -1, 64))
	}

	if r.SafeLevel != 0 {
		v.Add(options.KeySafeLevel, strconv.FormatFloat(r.SafeLevel, 'f', -1, 64))
	}

	return
}
