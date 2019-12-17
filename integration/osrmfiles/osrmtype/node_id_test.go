package osrmtype

import (
	"reflect"
	"testing"
)

func TestWriteNodeIDs(t *testing.T) {
	cases := []struct {
		p []byte
		n NodeIDs
	}{
		{
			[]byte{0x88, 0x83, 0x00, 0x00},
			NodeIDs{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, // redundant byte
			},
			NodeIDs{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			NodeIDs{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00, 0xd2, 0x14, 0x11, 0x00,
			},
			NodeIDs{33672, 1119442},
		},
	}

	for _, c := range cases {
		var n NodeIDs
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%nodeIDBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.n) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.n, n)
		}
	}
}
