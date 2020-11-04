package osmpatch

import (
	"testing"

	"github.com/qedus/osmpbf"
)

func TestIsValidWay(t *testing.T) {
	cases := []struct {
		way   *osmpbf.Way
		valid bool
	}{
		{nil, false},
		{&osmpbf.Way{}, false},
		{&osmpbf.Way{Tags: map[string]string{"highway": ""}}, true},
		{&osmpbf.Way{Tags: map[string]string{"route": ""}}, true},
		{&osmpbf.Way{Tags: map[string]string{"highway": "", "route": ""}}, true},
	}

	for _, c := range cases {
		b := IsValidWay(c.way)
		if b != c.valid {
			t.Errorf("way %v expect is valid %t but got %t", c.way, c.valid, b)
		}
	}
}
