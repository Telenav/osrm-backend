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
			[]Edge{Edge{From: 84760891102, To: 19496208102}},
			[]Edge{Edge{From: 19496208102, To: 84760891102}},
		},
		{
			[]Edge{
				Edge{From: 84762609102, To: 244183320001101},
				Edge{From: 244183320001101, To: 84762607102},
			},
			[]Edge{
				Edge{From: 84762607102, To: 244183320001101},
				Edge{From: 244183320001101, To: 84762609102},
			},
		},
		{
			[]Edge{
				Edge{From: 111, To: 84762609102},
				Edge{From: 84762609102, To: 244183320001101},
				Edge{From: 244183320001101, To: 84762607102},
			},
			[]Edge{
				Edge{From: 84762607102, To: 244183320001101},
				Edge{From: 244183320001101, To: 84762609102},
				Edge{From: 84762609102, To: 111},
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
