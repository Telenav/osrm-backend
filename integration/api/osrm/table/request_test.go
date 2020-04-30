package table

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/api/osrm/genericoptions"
)

func TestTableRequestURI(t *testing.T) {
	cases := []struct {
		r      Request
		expect string
	}{
		{
			Request{
				Service:     "table",
				Version:     "v1",
				Profile:     "driving",
				Coordinates: coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
		},
		{
			Request{
				Service:     "table",
				Version:     "v1",
				Profile:     "driving",
				Coordinates: coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769",
		},
		{
			Request{
				Service:     "table",
				Version:     "v1",
				Profile:     "driving",
				Coordinates: coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
				Sources:     genericoptions.Elements{"0"},
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769?sources=0",
		},
		{
			Request{
				Service:      "table",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
				Sources:      genericoptions.Elements{"0"},
				Destinations: genericoptions.Elements{"1", "2"},
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769?destinations=1;2&sources=0",
		},
		{
			Request{
				Service:      "table",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
				Sources:      genericoptions.Elements{"0"},
				Destinations: genericoptions.Elements{"1", "2"},
				Annotations:  "duration",
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769?annotations=duration&destinations=1;2&sources=0",
		},
		{
			Request{
				Service:      "table",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
				Sources:      genericoptions.Elements{"0"},
				Destinations: genericoptions.Elements{"1", "2"},
				Annotations:  "distance",
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769?annotations=distance&destinations=1;2&sources=0",
		},
		{
			Request{
				Service:      "table",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}, coordinate.Coordinate{Lat: 37.313769, Lon: -121.875655}},
				Sources:      genericoptions.Elements{"0"},
				Destinations: genericoptions.Elements{"1", "2"},
				Annotations:  "duration,distance",
			},
			"/table/v1/driving/-122.006349,37.364336;-121.875654,37.313767;-121.875655,37.313769?annotations=duration,distance&destinations=1;2&sources=0",
		},
	}

	for _, c := range cases {
		s := c.r.RequestURI()
		if s != c.expect {
			t.Errorf("%v QueryString(), expect %s, but got %s", c.r, c.expect, s)
		}
	}
}
