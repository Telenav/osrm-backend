package osrm

import (
	"reflect"
	"testing"
)

func TestPraseOptionElements(t *testing.T) {

	cases := []struct {
		s string
		OptionElements
		expectFail bool
	}{
		{"0;5;7", OptionElements{"0", "5", "7"}, false},
		{"0;5;7;", OptionElements{"0", "5", "7"}, false},
		{"0;;7", OptionElements{"0", "", "7"}, false},
		{"", OptionElements{}, false},
	}

	for _, c := range cases {
		elements, err := ParseOptionElemenets(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(elements, c.OptionElements) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.OptionElements, elements)
		}
	}
}

func TestOptionElementsString(t *testing.T) {
	cases := []struct {
		expect string
		OptionElements
	}{
		{"0;5;7", OptionElements{"0", "5", "7"}},
		{"0;;7", OptionElements{"0", "", "7"}},
		{"", OptionElements{""}},
		{"", OptionElements{}},
	}

	for _, c := range cases {
		s := c.OptionElements.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.OptionElements, c.expect, s)
		}
	}
}
