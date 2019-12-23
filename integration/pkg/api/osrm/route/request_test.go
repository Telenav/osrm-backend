package route

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route/options"
)

func TestRouteRequestURI(t *testing.T) {
	cases := []struct {
		r      Request
		expect string
	}{
		{
			Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesDefaultValue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
		},
		{
			Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesValueTrue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=true",
		},
		{
			Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: "100",
				Steps:        true,
				Annotations:  options.AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=100&steps=true",
		},
		{
			Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  options.AnnotationsValueTrue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=5&annotations=true&steps=true",
		},
	}

	for _, c := range cases {
		s := c.r.RequestURI()
		if s != c.expect {
			t.Errorf("%v QueryString(), expect %s, but got %s", c.r, c.expect, s)
		}
	}

}

func TestParseRouteRequest(t *testing.T) {

	cases := []struct {
		requestURI string
		expect     *Request
		expectFail bool
	}{
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?&alternatives=5&annotations=true&steps=true",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  options.AnnotationsValueTrue,
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=5&annotations=true&steps=true",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  options.AnnotationsValueTrue,
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesDefaultValue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesDefaultValue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			false,
		},
		{
			"route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesDefaultValue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=-1&annotations=tru,&steps=alse,",
			&Request{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  coordinate.Coordinates{coordinate.Coordinate{37.364336, -122.006349}, coordinate.Coordinate{37.313767, -121.875654}},
				Alternatives: options.AlternativesDefaultValue,
				Steps:        options.StepsDefaultValue,
				Annotations:  options.AnnotationsDefaultValue,
			},
			false,
		},
		{"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767.json", nil, true},
		{"/route/v1/driving/-122.006349,   37.364336;-121.875654,37.313767", nil, true},
		{"/route/v1/driving/-122.006349,37.364336;-121.875654,   37.313767?alternatives=-1", nil, true},
	}

	for _, c := range cases {
		r, err := ParseRequestURI(c.requestURI)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.requestURI, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(*r, *c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.requestURI, c.expect, r)
		}

	}
}
