package speedunit

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/util"
)

func TestConvert(t *testing.T) {

	cases := []struct {
		mps float64
		kph float64
	}{
		{0, 0},
		{1, 3.6},
		{30, 108},
	}

	for _, c := range cases {
		gotKPH := ConvertMPS2KPH(c.mps)
		if !util.Float64Equal(c.kph, gotKPH) {
			t.Errorf("ConvertMPS2KPH %f got %f but expect %f", c.mps, gotKPH, c.kph)
		}

		gotMPS := ConvertKPH2MPS(c.kph)
		if !util.Float64Equal(c.mps, gotMPS) {
			t.Errorf("ConvertKPH2MPS %f got %f but expect %f", c.kph, gotMPS, c.mps)
		}
	}
}
