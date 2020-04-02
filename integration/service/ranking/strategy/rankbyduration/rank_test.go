package rankbyduration

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

func TestRank(t *testing.T) {
	cases := []struct {
		input  []*route.Route
		expect []*route.Route
	}{
		{nil, nil},
		{
			[]*route.Route{&route.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil}},
			[]*route.Route{&route.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil}},
		},
		{
			[]*route.Route{
				&route.Route{Distance: 10000.0, Duration: 200.0, Geometry: "", Weight: 300.0, WeightName: "", Legs: nil},
				&route.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil},
				&route.Route{Distance: 22222.0, Duration: 300.0, Geometry: "", Weight: 200.0, WeightName: "", Legs: nil},
			},
			[]*route.Route{
				&route.Route{Distance: 12345.0, Duration: 100.0, Geometry: "", Weight: 100.0, WeightName: "", Legs: nil},
				&route.Route{Distance: 10000.0, Duration: 200.0, Geometry: "", Weight: 300.0, WeightName: "", Legs: nil},
				&route.Route{Distance: 22222.0, Duration: 300.0, Geometry: "", Weight: 200.0, WeightName: "", Legs: nil},
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
