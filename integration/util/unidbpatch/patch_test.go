package unidbpatch

import "testing"

func TestIsValidWay(t *testing.T) {
	cases := []struct {
		wayID int64
		valid bool
	}{
		{100, true},
		{106812895100, true},
		{-106812895100, true}, // should NOT exist in real case actually
		{0, false},
		{1000, false},
		{101, false},
		{106812895, false},
	}

	for _, c := range cases {
		b := IsValidWay(c.wayID)
		if b != c.valid {
			t.Errorf("wayID %d expect is valid %t but got %t", c.wayID, c.valid, b)
		}
	}
}

func TestTrimValidWayIDSuffix(t *testing.T) {
	cases := []struct {
		inWayID  int64
		outWayID int64
	}{
		{106812895100, 106812895},
		{-106812895100, -106812895}, // should NOT exist in real case actually
		{0, 0},
		{1000, 1000},
		{106812895, 106812895},
	}

	for _, c := range cases {
		wayID := TrimValidWayIDSuffix(c.inWayID)
		if wayID != c.outWayID {
			t.Errorf("trim suffix for wayID %d, expect %d but got %d", c.inWayID, c.outWayID, wayID)
		}
	}

}
