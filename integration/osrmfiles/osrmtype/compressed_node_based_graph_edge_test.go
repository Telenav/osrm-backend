package osrmtype

import (
	"reflect"
	"testing"
)

func TestCompressedNodeBasedGraphEdgesWrite(t *testing.T) {
	cases := []struct {
		p []byte
		CompressedNodeBasedGraphEdges
	}{
		{
			[]byte{0x0f, 0x00, 0x00, 0x00, 0x07, 0x0d, 0x0b, 0x00},
			CompressedNodeBasedGraphEdges{{Source: 15, Target: 724231}},
		},
		{
			[]byte{
				0x0f, 0x00, 0x00, 0x00, 0x07, 0x0d, 0x0b, 0x00,
				0xff, // redundant byte
			},
			CompressedNodeBasedGraphEdges{{Source: 15, Target: 724231}},
		},
		{
			[]byte{
				0x0f, 0x00, 0x00, 0x00, 0x07, 0x0d, 0x0b, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // redundant bytes
			},
			CompressedNodeBasedGraphEdges{{Source: 15, Target: 724231}},
		},
		{
			[]byte{
				0x0f, 0x00, 0x00, 0x00, 0x07, 0x0d, 0x0b, 0x00,
				0x0f, 0x00, 0x00, 0x00, 0xca, 0xfd, 0x08, 0x00,
			},
			CompressedNodeBasedGraphEdges{
				{Source: 15, Target: 724231},
				{Source: 15, Target: 589258},
			},
		},
	}

	for _, c := range cases {
		var edges CompressedNodeBasedGraphEdges
		_, err := edges.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(edges, c.CompressedNodeBasedGraphEdges) {
			t.Errorf("construct compressed node based graph edges from %v, expect %v, but got %v", c.p, c.CompressedNodeBasedGraphEdges, edges)
		}
	}
}
