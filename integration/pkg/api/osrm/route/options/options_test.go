package options

import (
	"reflect"
	"testing"
)

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
		alternatives, num, err := ParseAlternatives(c.s)
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

func TestParseSteps(t *testing.T) {
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
		b, err := ParseSteps(c.s)
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
		annotations, err := ParseAnnotations(c.s)
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
