package api

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	cases := []struct {
		s string
		url.Values
	}{
		{"waypoint=1", map[string][]string{"waypoint": []string{"1"}}},
		{"waypoint=1;2;3", map[string][]string{"waypoint": []string{"1;2;3"}}},
	}

	for _, c := range cases {
		v, _ := ParseQuery(c.s)
		if !reflect.DeepEqual(v, c.Values) {
			t.Errorf("parse %s, expect %v but got %v", c.s, c.Values, v)
		}
	}
}
