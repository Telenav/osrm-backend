package api

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseQuery is similar with url.ParseQuery(https://golang.org/pkg/net/url/#ParseQuery).
// The only difference is that url.ParseQuery uses both `&` and `;` as seperator,
// in contrast ParseQuery only uses `&`.
func ParseQuery(rawQuery string) (url.Values, error) {
	values := url.Values{}

	queryStr, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return values, err
	}

	for _, s := range strings.Split(queryStr, Ampersand) {
		keyValues := strings.Split(s, EqualTo)
		if len(keyValues) != 2 {
			if err == nil { // return err describes the first decoding error encountered, if any.
				err = fmt.Errorf("invalid query key-value %s", s)
			}
			continue
		}
		values[keyValues[0]] = []string{keyValues[1]}
	}
	return values, err
}

// ParseQueryDiscardError is similar with url.Query(https://golang.org/pkg/net/url/#URL.Query),
// which parses RawQuery and returns the corresponding values.
// It silently discards malformed value pairs. To check errors use ParseQuery.
func ParseQueryDiscardError(rawQuery string) url.Values {
	v, _ := ParseQuery(rawQuery)
	return v
}
