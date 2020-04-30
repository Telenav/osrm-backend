package nearbychargestation

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/Telenav/osrm-backend/integration/api/search/options"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
	"github.com/golang/glog"
)

// Request for search service
type Request struct {
	Service string
	Version string
	Action  string
	Format  string

	APIKey       string
	APISignature string
	Category     string
	Location     searchcoordinate.Coordinate
	Intent       string
	Locale       string
	Limit        int
	Radius       float64
}

// NewRequest create an empty entity Request.
func NewRequest() *Request {
	return &Request{
		// Path
		Service: "entity",
		Version: "v4",
		Action:  "search",
		Format:  "json",

		// Options
		APIKey:       "",
		APISignature: "",
		Category:     options.ChargeStationCategory,
		Location:     searchcoordinate.Coordinate{},
		Intent:       options.AroundIntent,
		Locale:       options.ENUSLocale,
		Limit:        options.DefaultLimitValue,
		Radius:       0.0,
	}
}

// RequestURI convert RouteRequest to RequestURI (e.g. "/path?foo=bar").
// see more in https://golang.org/pkg/net/url/#URL.RequestURI
func (r *Request) RequestURI() string {
	s := r.pathPrefix()

	queryStr := r.QueryString()
	if len(queryStr) > 0 {
		s += api.QuestionMark + queryStr
	}

	return s
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
		return nil, fmt.Errorf("empty URL input for nearbychargestaion->ParseRequestURL")
	}

	req := NewRequest()

	if err := req.parsePath(u.Path); err != nil {
		return nil, err
	}

	if values, err := url.ParseQuery(u.RawQuery); err == nil {
		req.parseQuery(values)
	} else {
		return nil, err
	}

	return req, nil
}

func (r *Request) pathPrefix() string {
	//i.e. "/entity/v4/search/json?"
	return api.Slash + r.Service + api.Slash + r.Version + api.Slash + r.Action + api.Slash + r.Format
}

// QueryString convert SearchRequest to "URL encoded" form ("bar=baz&foo=quux"), but NOT escape.
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
// i.e. "api_key=-&api_signature=-&query=-&location=-&intent=-&locale=-&limit=-&radius=-"
func (r *Request) QueryValues() (v url.Values) {
	v = make(url.Values)

	if len(r.APIKey) > 0 {
		v.Add(options.KeyAPIKey, r.APIKey)
	}

	if len(r.APISignature) > 0 {
		v.Add(options.KeyAPISignature, r.APISignature)
	}

	if len(r.Category) > 0 {
		v.Add(options.KeyQuery, r.Category)
	}

	coordinateStr := r.Location.String()
	if len(coordinateStr) > 0 {
		v.Add(options.KeyLocation, coordinateStr)
	}

	if len(r.Locale) > 0 {
		v.Add(options.KeyLocale, r.Locale)
	}

	if r.Limit > 0 {
		v.Add(options.KeyLimit, strconv.Itoa(r.Limit))
	}

	if r.Radius > 0 {
		v.Add(options.KeyRadius, strconv.FormatFloat(r.Radius, 'f', 6, 64))
	}

	return
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
	r.Action = s[2]
	r.Format = s[3]

	return nil
}

func (r *Request) parseQuery(values url.Values) {
	if v := values.Get(options.KeyAPIKey); len(v) > 0 {
		r.APIKey = v
	}

	if v := values.Get(options.KeyAPISignature); len(v) > 0 {
		r.APISignature = v
	}

	if v := values.Get(options.KeyLocation); len(v) > 0 {
		if location, err := searchcoordinate.ParseCoordinate(v); err == nil {
			r.Location = location
		}
	}

	if v := values.Get(options.KeyLocale); len(v) > 0 {
		r.Locale = v
	}

	if v := values.Get(options.KeyIntent); len(v) > 0 {
		r.Intent = v
	}

	if v := values.Get(options.KeyLimit); len(v) > 0 {
		if limit, err := strconv.Atoi(v); err == nil {
			r.Limit = limit
		}
	}

	if v := values.Get(options.KeyRadius); len(v) > 0 {
		if radius, err := strconv.ParseFloat(v, 64); err == nil {
			r.Radius = radius
		}
	}

}
