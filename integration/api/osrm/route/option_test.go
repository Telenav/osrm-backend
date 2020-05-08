package route

import (
	"reflect"
	"testing"
)

func TestParseOptionAlternatives(t *testing.T) {

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
		alternatives, num, err := parseOptionAlternatives(c.s)
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

func TestParseOptionSteps(t *testing.T) {
	cases := []struct {
		s          string
		expect     bool
		expectFail bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"", false, true},
		{"true1", false, true},
		{"-1", false, true},
	}

	for _, c := range cases {
		b, err := parseOptionSteps(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if b != c.expect {
			t.Errorf("parse %s, expect %t, but got %t", c.s, c.expect, b)
		}
	}
}

func TestParseOptionAnnotations(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectFail bool
	}{
		{"true", OptionValueTrue, false},
		{"false", OptionValueFalse, false},
		{"nodes", "nodes", false},
		{"nodes,distance", "nodes,distance", false},
		{"nodes,distance,duration,datasources,weight,speed", "nodes,distance,duration,datasources,weight,speed", false},
		{"nodes,", "", true},
		{"nodes,distances", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		annotations, err := parseOptionAnnotations(c.s)
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

func TestParseOptionGeometries(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectFail bool
	}{
		{"polyline", OptionGeometriesValuePolyline, false},
		{"polyline6", OptionGeometriesValuePolyline6, false},
		{"geojson", OptionGeometriesValueGeojson, false},
		{"polyli", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		geometries, err := parseOptionGeometries(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(geometries, c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.expect, geometries)
		}
	}
}

func TestParseOptionOverview(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectFail bool
	}{
		{"simplified", OptionOverviewValueSimplified, false},
		{"full", OptionOverviewValueFull, false},
		{"false", OptionOverviewValueFalse, false},
		{"simp", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		overview, err := parseOptionOverview(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(overview, c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.expect, overview)
		}
	}
}

func TestParseOptionContinueStraight(t *testing.T) {

	cases := []struct {
		s          string
		expect     string
		expectFail bool
	}{
		{"default", OptionContinueStraightValueDefault, false},
		{"false", OptionContinueStraightValueFalse, false},
		{"true", OptionContinueStraightValueTrue, false},
		{"1", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		continueStraight, err := parseOptionContinueStraight(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(continueStraight, c.expect) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.expect, continueStraight)
		}
	}
}
