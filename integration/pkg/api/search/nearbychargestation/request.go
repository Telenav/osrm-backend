package nearbychargestation

import (
	"net/url"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/options"
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
	Location     coordinate.Coordinate
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
		Location:     coordinate.Coordinate{},
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
