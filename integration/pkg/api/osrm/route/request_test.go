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
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesValueTrue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=true",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesValuePolyline6,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?geometries=polyline6",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewValueFull,
				ContinueStraight: options.ContinueStraightDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?overview=full",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightValueFalse,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?continue_straight=false",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
				Waypoints:        coordinate.Indexes{0, 3, 5},
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?waypoints=0;3;5",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     "100",
				Steps:            true,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=100&steps=true",
		},
		{
			Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     "5",
				Steps:            true,
				Annotations:      options.AnnotationsValueTrue,
				Geometries:       options.GeometriesValueGeojson,
				Overview:         options.OverviewValueFalse,
				ContinueStraight: options.ContinueStraightValueTrue,
				Waypoints:        coordinate.Indexes{2, 3},
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=5&annotations=true&continue_straight=true&geometries=geojson&overview=false&steps=true&waypoints=2;3",
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
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?&alternatives=5&annotations=true&steps=true&geometries=polyline&overview=false&continue_straight=true&waypoints=2;4;6",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     "5",
				Steps:            true,
				Annotations:      options.AnnotationsValueTrue,
				Geometries:       options.GeometriesValuePolyline,
				Overview:         options.OverviewValueFalse,
				ContinueStraight: options.ContinueStraightValueTrue,
				Waypoints:        coordinate.Indexes{2, 4, 6},
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=5&annotations=true&steps=true&geometries=polyline&overview=false&continue_straight=true&waypoints=2;4;6",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     "5",
				Steps:            true,
				Annotations:      options.AnnotationsValueTrue,
				Geometries:       options.GeometriesValuePolyline,
				Overview:         options.OverviewValueFalse,
				ContinueStraight: options.ContinueStraightValueTrue,
				Waypoints:        coordinate.Indexes{2, 4, 6},
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
				Waypoints:        coordinate.Indexes{},
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
				Waypoints:        coordinate.Indexes{},
			},
			false,
		},
		{
			"route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
				Waypoints:        coordinate.Indexes{},
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=-1&annotations=tru,&steps=alse,",
			&Request{
				Service:          "route",
				Version:          "v1",
				Profile:          "driving",
				Coordinates:      coordinate.Coordinates{coordinate.Coordinate{Lat: 37.364336, Lon: -122.006349}, coordinate.Coordinate{Lat: 37.313767, Lon: -121.875654}},
				Alternatives:     options.AlternativesDefaultValue,
				Steps:            options.StepsDefaultValue,
				Annotations:      options.AnnotationsDefaultValue,
				Geometries:       options.GeometriesDefaultValue,
				Overview:         options.OverviewDefaultValue,
				ContinueStraight: options.ContinueStraightDefaultValue,
				Waypoints:        coordinate.Indexes{},
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
