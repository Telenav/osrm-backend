package osrm

import (
	"reflect"
	"testing"
)

func TestPraseOptionClasses(t *testing.T) {

	cases := []struct {
		s string
		OptionClasses
		expectFail bool
	}{
		{"0,5,7", OptionClasses{"0", "5", "7"}, false},
		{"0,5,7,", OptionClasses{"0", "5", "7"}, false},
		{"0,,7", OptionClasses{"0", "7"}, false},
		{"", OptionClasses{}, false},
	}

	for _, c := range cases {
		classes, err := ParseOptionClasses(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(classes, c.OptionClasses) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.OptionClasses, classes)
		}
	}
}

func TestOptionClassesString(t *testing.T) {
	cases := []struct {
		expect string
		OptionClasses
	}{
		{"0,5,7", OptionClasses{"0", "5", "7"}},
		{"0,7", OptionClasses{"0", "", "7"}},
		{"", OptionClasses{""}},
		{"", OptionClasses{}},
	}

	for _, c := range cases {
		s := c.OptionClasses.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.OptionClasses, c.expect, s)
		}
	}
}
