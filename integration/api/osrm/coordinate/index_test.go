package coordinate

import (
	"reflect"
	"testing"
)

func TestPraseIndexes(t *testing.T) {

	cases := []struct {
		s string
		Indexes
		expectFail bool
	}{
		{"0;5;7", Indexes{0, 5, 7}, false},
		{"0;5;7;", Indexes{0, 5, 7}, false},
		{"0", Indexes{0}, false},
		{"5;1;4;2;3;6", Indexes{5, 1, 4, 2, 3, 6}, false},
		{"", Indexes{}, false},
		{"-1", nil, true},
		{"a", nil, true},
	}

	for _, c := range cases {
		indexes, err := PraseIndexes(c.s)
		if err != nil && c.expectFail {
			continue //right
		} else if (err != nil && !c.expectFail) || (err == nil && c.expectFail) {
			t.Errorf("parse %s expect fail %t, but got err %v", c.s, c.expectFail, err)
			continue
		}

		if !reflect.DeepEqual(indexes, c.Indexes) {
			t.Errorf("parse %s, expect %v, but got %v", c.s, c.Indexes, indexes)
		}
	}
}

func TestIndexesString(t *testing.T) {
	cases := []struct {
		expect string
		Indexes
	}{
		{"0;5;7", Indexes{0, 5, 7}},
		{"0", Indexes{0}},
		{"5;1;4;2;3;6", Indexes{5, 1, 4, 2, 3, 6}},
		{"", Indexes{}},
	}

	for _, c := range cases {
		s := c.Indexes.String()
		if s != c.expect {
			t.Errorf("%v String(), expect %s, but got %s", c.Indexes, c.expect, s)
		}
	}
}
