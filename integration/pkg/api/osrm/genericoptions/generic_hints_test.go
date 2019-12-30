package genericoptions

import "testing"

func TestParseGenerateHints(t *testing.T) {
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
		b, err := ParseGenerateHints(c.s)
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
