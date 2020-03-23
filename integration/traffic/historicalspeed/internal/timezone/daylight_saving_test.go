package timezone

import "testing"

func TestParseDaylightSaving(t *testing.T) {
	cases := []struct {
		s     string
		dst   int8
		valid bool
	}{
		{"0", 0, true},
		{"67", 67, true},
		{"-1", 0, false},
		{"68", 0, false},
	}

	for _, c := range cases {
		dst, err := ParseDaylightSaving(c.s)
		if err != nil {
			if c.valid {
				t.Errorf("parse daylight saving %s, expect succeed but got err %v", c.s, err)
			}
			continue
		}

		if !c.valid {
			t.Errorf("parse daylight saving %s, expect failed but succeed", c.s)
		}

		if dst != c.dst {
			t.Errorf("parse daylight saving %s, expect %d but got %d", c.s, c.dst, dst)
		}
	}
}

func TestFormatDaylightSaving(t *testing.T) {

	cases := []struct {
		dst int8
		s   string
	}{
		{0, "0"},
		{67, "67"},
	}

	for _, c := range cases {
		s := FormatDaylightSaving(c.dst)
		if s != c.s {
			t.Errorf("format daylight saving %d, expect %s but got %s", c.dst, c.s, s)
		}
	}
}

func TestIsValidDaylightSaving(t *testing.T) {

	cases := []struct {
		dst     int8
		isValid bool
	}{
		{0, true},
		{67, true},
		{-1, false},
		{68, false},
	}

	for _, c := range cases {
		isValid := isValidDaylightSaving(c.dst)
		if isValid != c.isValid {
			t.Errorf("vaidate daylightSaving %d, expect %t but got %t", c.dst, c.isValid, isValid)
		}
	}
}
