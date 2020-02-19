package oasis

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
)

var mockFloatArray1 []float64 = []float64{22.2, 11.1, 33.3, 44.4}
var mock1ToNTableResponse1 table.Response = table.Response{
	Durations: [][]*float64{
		[]*float64{&mockFloatArray1[0], &mockFloatArray1[1], &mockFloatArray1[2], &mockFloatArray1[3]},
	},
	Distances: [][]*float64{
		[]*float64{&mockFloatArray1[0], &mockFloatArray1[1], &mockFloatArray1[2], &mockFloatArray1[3]},
	},
}

var mockFloatArray2 []float64 = []float64{66.6, 11.1, 33.3, 33.3}
var mockNTo1TableResponse1 table.Response = table.Response{
	Durations: [][]*float64{
		[]*float64{&mockFloatArray2[0]},
		[]*float64{&mockFloatArray2[1]},
		[]*float64{&mockFloatArray2[2]},
		[]*float64{&mockFloatArray2[3]},
	},
	Distances: [][]*float64{
		[]*float64{&mockFloatArray2[0]},
		[]*float64{&mockFloatArray2[1]},
		[]*float64{&mockFloatArray2[2]},
		[]*float64{&mockFloatArray2[3]},
	},
}

func TestRankingSingleChargeStation(t *testing.T) {
	index, err := rankingSingleChargeStation(&mock1ToNTableResponse1, &mockNTo1TableResponse1)
	if err != nil || index != 1 {
		t.Errorf("expect %v but got %v", 1, index)
	}
}
