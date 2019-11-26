package rankbyduration

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
)

func TestRank(t *testing.T) {
	cases := []struct {
		input  []*osrmv1.Route
		expect []*osrmv1.Route
	}{
		{nil, nil},
		{
			[]*osrmv1.Route{&osrmv1.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil}},
			[]*osrmv1.Route{&osrmv1.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil}},
		},
		{
			[]*osrmv1.Route{
				&osrmv1.Route{Distance: 10000.0, Duration: 200.0, Geometry: "", Weight: 300.0, WeightName: "", Legs: nil},
				&osrmv1.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil},
				&osrmv1.Route{Distance: 22222.0, Duration: 300.0, Geometry: "", Weight: 200.0, WeightName: "", Legs: nil},
			},
			[]*osrmv1.Route{
				&osrmv1.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil},
				&osrmv1.Route{Distance: 10000.0, Duration: 200.0, Geometry: "", Weight: 300.0, WeightName: "", Legs: nil},
				&osrmv1.Route{Distance: 22222.0, Duration: 300.0, Geometry: "", Weight: 200.0, WeightName: "", Legs: nil},
			},
		},
	}

	for _, c := range cases {
		r := Rank(c.input)
		if !reflect.DeepEqual(c.expect, r) {
			t.Errorf("rank %v by duration expect %v, but got %v", c.input, c.expect, r)
		}
	}
}
