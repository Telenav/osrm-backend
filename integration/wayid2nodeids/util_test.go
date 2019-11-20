package wayid2nodeids

import "testing"

func TestAbsInt64(t *testing.T) {

	cases := []struct {
		in     int64
		expect int64
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
	}

	for _, c := range cases {
		out := absInt64(c.in)
		if out != c.expect {
			t.Errorf("expect %d for %d after absInt64, but got %d", c.expect, c.in, out)
		}
	}
}
