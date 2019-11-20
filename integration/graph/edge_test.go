package graph

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	cases := []struct {
		in     []Edge
		expect []Edge
	}{
		{nil, nil},
		{[]Edge{}, []Edge{}},
		{
			[]Edge{Edge{FromNode: 84760891102, ToNode: 19496208102}},
			[]Edge{Edge{FromNode: 19496208102, ToNode: 84760891102}},
		},
		{
			[]Edge{
				Edge{FromNode: 84762609102, ToNode: 244183320001101},
				Edge{FromNode: 244183320001101, ToNode: 84762607102},
			},
			[]Edge{
				Edge{FromNode: 84762607102, ToNode: 244183320001101},
				Edge{FromNode: 244183320001101, ToNode: 84762609102},
			},
		},
		{
			[]Edge{
				Edge{FromNode: 111, ToNode: 84762609102},
				Edge{FromNode: 84762609102, ToNode: 244183320001101},
				Edge{FromNode: 244183320001101, ToNode: 84762607102},
			},
			[]Edge{
				Edge{FromNode: 84762607102, ToNode: 244183320001101},
				Edge{FromNode: 244183320001101, ToNode: 84762609102},
				Edge{FromNode: 84762609102, ToNode: 111},
			},
		},
	}

	for _, c := range cases {
		out := ReverseEdges(c.in)
		if !reflect.DeepEqual(out, c.expect) {
			t.Errorf("expect %v for %v after reverse, but got %v", c.expect, c.in, out)
		}
	}

}
