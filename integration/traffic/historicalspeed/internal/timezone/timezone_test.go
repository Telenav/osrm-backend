package timezone

import "testing"

func TestParseTimezone(t *testing.T) {
	cases := []struct {
		s     string
		tz    int16
		valid bool
	}{
		{"-120", -120, true},
		{"-70", -70, true},
		{"-070", -70, true},
		{"-035", -35, true},
		{"000", 0, true},
		{"0", 0, true},
		{"080", 80, true},
		{"80", 80, true},
		{"127", 127, true},
		{"140", 140, true},
		{"-125", 0, false},
		{"145", 0, false},
	}

	for _, c := range cases {
		tz, err := ParseTimezone(c.s)
		if err != nil {
			if c.valid {
				t.Errorf("parse timezone %s, expect succeed but got err %v", c.s, err)
			}
			continue
		}

		if !c.valid {
			t.Errorf("parse timezone %s, expect failed but succeed", c.s)
		}

		if tz != c.tz {
			t.Errorf("parse timezone %s, expect %d but got %d", c.s, c.tz, tz)
		}
	}
}

func TestFormatTimezone(t *testing.T) {

	cases := []struct {
		tz int16
		s  string
	}{
		{-120, "-120"},
		{-70, "-070"},
		{0, "000"},
		{80, "080"},
		{140, "140"},
	}

	for _, c := range cases {
		s := FormatTimezone(c.tz)
		if s != c.s {
			t.Errorf("format timezone %d, expect %s but got %s", c.tz, c.s, s)
		}
	}
}

func TestIsValidTimezone(t *testing.T) {

	cases := []struct {
		tz      int16
		isValid bool
	}{
		{-120, true},
		{-70, true},
		{-35, true},
		{0, true},
		{80, true},
		{127, true},
		{140, true},
		{-125, false},
		{145, false},
	}

	for _, c := range cases {
		isValid := isValidTimezone(c.tz)
		if isValid != c.isValid {

			//better format timezone
			if c.tz >= 0 {
				t.Errorf("vaidate timezone %03d, expect %t but got %t", c.tz, c.isValid, isValid)
			} else {
				t.Errorf("vaidate timezone %04d, expect %t but got %t", c.tz, c.isValid, isValid)
			}
		}
	}
}
