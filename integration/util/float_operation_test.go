package util

import "testing"

func TestFloat64Equal(t *testing.T) {
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
		actual := Float64Equal(c.num1, c.num2)
		if actual != c.expect {
			t.Errorf("TestFloat64Equal failed with case \n %#v expect \n %#v but got %#v\n", c, c.expect, actual)
		}
	}

}

func TestFloat32Equal(t *testing.T) {
	cases := []struct {
		num1   float32
		num2   float32
		expect bool
	}{
		{
			32.33333,
			32.33333,
			true,
		},
		{
			32.3332,
			32.3333,
			false,
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
			-32.3333333,
			-32.3333334,
			true,
		},
	}

	for _, c := range cases {
		actual := Float32Equal(c.num1, c.num2)
		if actual != c.expect {
			t.Errorf("TestFloat64Equal failed with case \n %#v expect \n %#v but got %#v\n", c, c.expect, actual)
		}
	}

}
