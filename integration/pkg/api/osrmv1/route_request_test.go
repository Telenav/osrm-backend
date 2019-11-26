package osrmv1

import (
	"reflect"
	"testing"
)

func TestRouteRequestURI(t *testing.T) {
	cases := []struct {
		r      RouteRequest
		expect string
	}{
		{
			RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesDefaultValue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
		},
		{
			RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesValueTrue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=true",
		},
		{
			RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: "100",
				Steps:        true,
				Annotations:  AnnotationsDefaultValue,
			},
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=100&steps=true",
		},
		{
			RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  AnnotationsValueTrue,
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
		expect     *RouteRequest
		expectFail bool
	}{
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?&alternatives=5&annotations=true&steps=true",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  AnnotationsValueTrue,
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=5&annotations=true&steps=true",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: "5",
				Steps:        true,
				Annotations:  AnnotationsValueTrue,
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesDefaultValue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			false,
		},
		{
			"http://localhost:8080/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesDefaultValue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			false,
		},
		{
			"route/v1/driving/-122.006349,37.364336;-121.875654,37.313767",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesDefaultValue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			false,
		},
		{
			"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767?alternatives=-1&annotations=tru,&steps=alse,",
			&RouteRequest{
				Service:      "route",
				Version:      "v1",
				Profile:      "driving",
				Coordinates:  Coordinates{Coordinate{37.364336, -122.006349}, Coordinate{37.313767, -121.875654}},
				Alternatives: AlternativesDefaultValue,
				Steps:        StepsDefaultValue,
				Annotations:  AnnotationsDefaultValue,
			},
			false,
		},
		{"/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767.json", nil, true},
		{"/route/v1/driving/-122.006349,   37.364336;-121.875654,37.313767", nil, true},
		{"/route/v1/driving/-122.006349,37.364336;-121.875654,   37.313767?alternatives=-1", nil, true},
	}

	for _, c := range cases {
		r, err := ParseRouteRequestURI(c.requestURI)
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

func TestParseAlternatives(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectNum  int
		expectFail bool
	}{
		{"true", "true", 2, false},
		{"false", "false", 1, false},
		{"0", "0", 0, false}, // same as false
		{"1", "1", 1, false},
		{"5", "5", 5, false},
		{"100", "100", 100, false},
		{"111111111", "111111111", 111111111, false},
		{"-1", "", 1, true},
	}

	for _, c := range cases {
		alternatives, num, err := parseAlternatives(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if alternatives != c.expect || num != c.expectNum {
			t.Errorf("parse %s, expect %s %d, but got %s %d", c.s, c.expect, c.expectNum, alternatives, num)
		}
	}
}

func TestParseAnnotations(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectFail bool
	}{
		{"true", ValueTrue, false},
		{"false", ValueFalse, false},
		{"nodes", "nodes", false},
		{"nodes,distance", "nodes,distance", false},
		{"nodes,distance,duration,datasources,weight,speed", "nodes,distance,duration,datasources,weight,speed", false},
		{"nodes,", "", true},
		{"nodes,distances", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		annotations, err := parseAnnotations(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(annotations, c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.expect, annotations)
		}
	}

}
