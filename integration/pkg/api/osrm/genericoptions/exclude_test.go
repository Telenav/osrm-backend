package genericoptions

import (
	"reflect"
	"testing"
)

func TestPraseClasses(t *testing.T) {

	cases := []struct {
		s string
		Classes
		expectFail bool
	}{
		{"0,5,7", Classes{"0", "5", "7"}, false},
		{"0,5,7,", Classes{"0", "5", "7"}, false},
		{"0,,7", Classes{"0", "7"}, false},
		{"", Classes{}, false},
	}

	for _, c := range cases {
		classes, err := ParseClasses(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(classes, c.Classes) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.Classes, classes)
		}
	}
}

func TestClassesString(t *testing.T) {
	cases := []struct {
		expect string
		Classes
	}{
		{"0,5,7", Classes{"0", "5", "7"}},
		{"0,7", Classes{"0", "", "7"}},
		{"", Classes{""}},
		{"", Classes{}},
	}

	for _, c := range cases {
		s := c.Classes.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.Classes, c.expect, s)
		}
	}
}
