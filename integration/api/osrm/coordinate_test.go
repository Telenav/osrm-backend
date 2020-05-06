package osrm

import (
	"reflect"
	"testing"
)

func TestCoordinatesString(t *testing.T) {
	cases := []struct {
		s string
		c Coordinates
	}{
		{"13.388860,52.517037", Coordinates{Coordinate{52.517037, 13.388860}}},
		{
			"13.388860,52.517037;13.397634,52.529407;13.428555,52.523219",
			Coordinates{
				Coordinate{52.517037, 13.388860},
				Coordinate{52.529407, 13.397634},
				Coordinate{52.523219, 13.428555},
			},
		},
	}

	for _, c := range cases {
		s := c.c.String()
		if s != c.s {
			t.Errorf("%v String(), expect %s, but got %s", c.c, c.s, s)
		}
	}
}

func TestParseCoordinates(t *testing.T) {
	cases := []struct {
		s          string
		expect     Coordinates
		expectFail bool
	}{
		{"", Coordinates{}, true},
		{"13.388860", Coordinates{}, true},
		{"13.388860;52.517037", Coordinates{}, true},
		{"a13.388860,52.517037", Coordinates{}, true},
		{"13.388860,52.517037a", Coordinates{}, true},
		{"13.388860,52.517037.", Coordinates{}, true},

		{"13.388860,52.517037", Coordinates{Coordinate{52.517037, 13.388860}}, false},
		{"13.388860,52.517037", Coordinates{Coordinate{52.517037, 13.388860}}, false},
		{
			"13.388860,52.517037;13.397634,52.529407;13.428555,52.523219",
			Coordinates{
				Coordinate{52.517037, 13.388860},
				Coordinate{52.529407, 13.397634},
				Coordinate{52.523219, 13.428555},
			},
			false,
		},
		{
			"13.388860,52.517037;13.397634,52.529407;13.428555,52.523219.json",
			Coordinates{},
			true,
		},
	}

	for _, c := range cases {
		coordinates, err := ParseCoordinates(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(coordinates, c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.expect, coordinates)
		}
	}

}

func TestPraseCoordinateIndexes(t *testing.T) {

	cases := []struct {
		s string
		CoordinateIndexes
		expectFail bool
	}{
		{"0;5;7", CoordinateIndexes{0, 5, 7}, false},
		{"0;5;7;", CoordinateIndexes{0, 5, 7}, false},
		{"0", CoordinateIndexes{0}, false},
		{"5;1;4;2;3;6", CoordinateIndexes{5, 1, 4, 2, 3, 6}, false},
		{"", CoordinateIndexes{}, false},
		{"-1", nil, true},
		{"a", nil, true},
	}

	for _, c := range cases {
		indexes, err := PraseCoordinateIndexes(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(indexes, c.CoordinateIndexes) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.CoordinateIndexes, indexes)
		}
	}
}

func TestCoordinateIndexesString(t *testing.T) {
	cases := []struct {
		expect string
		CoordinateIndexes
	}{
		{"0;5;7", CoordinateIndexes{0, 5, 7}},
		{"0", CoordinateIndexes{0}},
		{"5;1;4;2;3;6", CoordinateIndexes{5, 1, 4, 2, 3, 6}},
		{"", CoordinateIndexes{}},
	}

	for _, c := range cases {
		s := c.CoordinateIndexes.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.CoordinateIndexes, c.expect, s)
		}
	}
}
