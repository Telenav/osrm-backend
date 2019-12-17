package osrmtype

import (
	"math"
	"reflect"
	"testing"
)

func TestNodeBasedEdgesWrite(t *testing.T) {

	cases := []struct {
		p []byte
		NodeBasedEdges
	}{
		{
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x84, 0x70, 0x07, 0x00, 0x18, 0x00, 0x00, 0x00, 0x18, 0x00, 0x00, 0x00,
				0x73, 0x25, 0x5b, 0x42, 0xff, 0xff, 0xff, 0x7f, 0xce, 0xad, 0x00, 0x00, 0xa1, 0x01, 0x02, 0x00,
			},
			NodeBasedEdges{
				NodeBasedEdge{
					Source:         0,
					Target:         487556,
					Weight:         24,
					Duration:       24,
					Distance:       54.7865715,
					GeometryID:     GeometryID{math.MaxUint32 >> 1, false},
					AnnotationData: 44494,
					//TODO: NodeBasedEdgeClassification
				},
			},
		},
	}

	for _, c := range cases {
		var n NodeBasedEdges
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%nodeBasedEdgeBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.NodeBasedEdges) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.NodeBasedEdges, n)
		}
	}
}
