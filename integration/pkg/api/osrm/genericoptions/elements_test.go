package genericoptions

import (
	"reflect"
	"testing"
)

func TestPraseElements(t *testing.T) {

	cases := []struct {
		s string
		Elements
		expectFail bool
	}{
		{"0;5;7", Elements{"0", "5", "7"}, false},
		{"0;5;7;", Elements{"0", "5", "7"}, false},
		{"0;;7", Elements{"0", "", "7"}, false},
		{"", Elements{}, false},
	}

	for _, c := range cases {
		elements, err := ParseElemenets(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(elements, c.Elements) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.Elements, elements)
		}
	}
}

func TestElementsString(t *testing.T) {
	cases := []struct {
		expect string
		Elements
	}{
		{"0;5;7", Elements{"0", "5", "7"}},
		{"0;;7", Elements{"0", "", "7"}},
		{"", Elements{""}},
		{"", Elements{}},
	}

	for _, c := range cases {
		s := c.Elements.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.Elements, c.expect, s)
		}
	}
}
