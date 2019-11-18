package wayid2nodeids

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {

	cases := []struct {
		line string
		ids  []int64
	}{
		{"", nil},
		{"24418325,84760891102", nil},
		{"24418325,84760891102,19496208102", []int64{24418325, 84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,", []int64{24418325, 84760891102, 19496208102}},
		{"24418325,84760891102,19496208102,,,,,", []int64{24418325, 84760891102, 19496208102}},
	}

	for _, c := range cases {
		result := parseLine(c.line)
		if !reflect.DeepEqual(result, c.ids) {
			t.Errorf("parseLine %s, expect %v, but got %v", c.line, c.ids, result)
		}
	}
}
