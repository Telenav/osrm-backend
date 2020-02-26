package nametable

import "testing"

// unittest in C++ implementation: https://github.com/Telenav/osrm-backend/blob/b24b8a085dc10bea279ffb352049330beae23791/unit_tests/util/indexed_data.cpp#L22
func TestSum2Bits(t *testing.T) {
	cases := []struct {
		v   uint32
		sum uint32
	}{
		{0xe4, 6},
		{0x11111111, 8},
		{0x55555555, 16},
		{0xffffffff, 48},
		{0x55556AAB, 24},
	}

	for _, c := range cases {
		sum := sum2Bits(c.v)
		if sum != c.sum {
			t.Errorf("sum2Bits(%d)=%d, but want %d", c.v, sum, c.sum)
		}
	}
}
