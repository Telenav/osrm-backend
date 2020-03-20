package historicalspeed

import "testing"

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
