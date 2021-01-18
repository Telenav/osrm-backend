package osrmtype

import (
	"reflect"
	"testing"
)

func TestEdgeBasedNodesWrite(t *testing.T) {

	cases := []struct {
		p []byte
		EdgeBasedNodes
	}{
		{
			[]byte{
				0x00, 0x00, 0x00, 0x80, 0x08, 0x01, 0x00, 0x00, 0x2F, 0x00, 0x00, // insufficient data
			},
			nil,
		},
		{
			[]byte{
				0x00, 0x00, 0x00, 0x80, 0x08, 0x01, 0x00, 0x00, 0x2F, 0x00, 0x00, 0x00,
			},
			EdgeBasedNodes{
				EdgeBasedNode{
					GeometryID:   GeometryID{0, true},
					ComponentID:  ComponentID{264, false},
					AnnotationID: 47,
					Segregated:   false,
				},
			},
		},
		{
			[]byte{
				0x00, 0x00, 0x00, 0x80, 0x08, 0x01, 0x00, 0x00, 0x2F, 0x00, 0x00, 0x80,
			},
			EdgeBasedNodes{
				EdgeBasedNode{
					GeometryID:   GeometryID{0, true},
					ComponentID:  ComponentID{264, false},
					AnnotationID: 47,
					Segregated:   true,
				},
			},
		},
		{
			[]byte{
				0x00, 0x00, 0x00, 0x80, 0x08, 0x01, 0x00, 0x00, 0x2F, 0x00, 0x00, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, // redundant bytes
			},
			EdgeBasedNodes{
				EdgeBasedNode{
					GeometryID:   GeometryID{0, true},
					ComponentID:  ComponentID{264, false},
					AnnotationID: 47,
					Segregated:   false,
				},
			},
		},
	}

	for _, c := range cases {
		var n EdgeBasedNodes
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%edgeBasedNodeBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.EdgeBasedNodes) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.EdgeBasedNodes, n)
		}
	}
}
