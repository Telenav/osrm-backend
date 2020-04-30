package nearbychargestation

import (
	"testing"
)

func TestNearBySearchURIDecoding(t *testing.T) {
	cases := []string{
		"/entity/v4/search/json?api_key=key&api_signature=signature&limit=25&locale=en-US&location=37.785090,-122.419880&query=EV+Charging+Station",
	}

	for _, c := range cases {
		req, _ := ParseRequestURI(c)
		if c != req.RequestURI() {
			t.Errorf("During test ParseRequestURI(), expect %s, but got %s", c, req.RequestURI())
		}
	}
}
