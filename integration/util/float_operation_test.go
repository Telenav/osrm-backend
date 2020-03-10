package util

import "testing"

func TestFloatEquals(t *testing.T) {
	cases := []struct {
		num1   float64
		num2   float64
		expect bool
	}{
		{
			32.33333,
			32.33333,
			true,
		},
		{
			32.33333333,
			32.33333334,
			false,
		},
		{
			32.333333333,
			32.333333334,
			true,
		},
		{
			32.33333,
			-32.33333,
			false,
		},
		{
			-32.33333,
			-32.33333,
			true,
		},
		{
			-32.33333333,
			-32.33333334,
			false,
		},
		{
			-32.333333333,
			-32.333333334,
			true,
		},
	}

	for _, c := range cases {
		actual := FloatEquals(c.num1, c.num2)
		if actual != c.expect {
			t.Errorf("TestFloatEquals failed with case \n %#v expect \n %#v but got %#v\n", c, c.expect, actual)
		}
	}

}
