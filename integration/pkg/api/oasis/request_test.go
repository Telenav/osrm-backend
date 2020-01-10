package oasis

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
)

func TestOasisRequestURI(t *testing.T) {
	cases := []struct {
		r      Request
		expect string
	}{
		{
			Request{
				Service:     "oasis",
				Version:     "v1",
				Profile:     "earliest",
				Coordinates: coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
			},
			"/oasis/v1/earliest/-122.006349,37.364336;-121.875654,37.313767",
		},
		{
			Request{
				Service:     "oasis",
				Version:     "v1",
				Profile:     "earliest",
				Coordinates: coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				MaxRange:    300.1,
				CurrRange:   100,
				PreferLevel: 80.0,
				SafeLevel:   50.0,
			},
			"/oasis/v1/earliest/-122.006349,37.364336;-121.875654,37.313767?curr_range=100&max_range=300.1&prefer_level=80&safe_level=50",
		},
	}

	for _, c := range cases {
		s := c.r.RequestURI()
		if s != c.expect {
			t.Errorf("%v QueryString(), expect %s, but got %s", c.r, c.expect, s)
		}
	}
}
