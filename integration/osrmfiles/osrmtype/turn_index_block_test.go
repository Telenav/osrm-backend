package osrmtype

import (
	"reflect"
	"testing"
)

func TestWriteTurnIndexBlocks(t *testing.T) {
	cases := []struct {
		p []byte
		n TurnIndexBlocks
	}{
		{
			[]byte{0x6F, 0xAC, 0x03, 0x00, 0x10, 0x95, 0x12, 0x00, 0xC1, 0xA0, 0x03, 0x00},
			TurnIndexBlocks{[]TurnIndexBlock{{From: 240751, Via: 1217808, To: 237761}}, []byte{}},
		},
		{
			[]byte{
				0x6F, 0xAC, 0x03, 0x00, 0x10, 0x95, 0x12, 0x00, 0xC1, 0xA0, 0x03, 0x00,
				0xFF, // redundant byte
			},
			TurnIndexBlocks{[]TurnIndexBlock{{From: 240751, Via: 1217808, To: 237761}}, []byte{0xFF}},
		},
		{
			[]byte{
				0x6F, 0xAC, 0x03, 0x00, 0x10, 0x95, 0x12, 0x00, 0xC1, 0xA0, 0x03, 0x00,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			TurnIndexBlocks{[]TurnIndexBlock{{From: 240751, Via: 1217808, To: 237761}}, []byte{0xFF, 0xFF, 0xFF}},
		},
		{
			[]byte{
				0x6F, 0xAC, 0x03, 0x00, 0x10, 0x95, 0x12, 0x00, 0xC1, 0xA0, 0x03, 0x00, 0x6F, 0xAC, 0x03, 0x00, 0x10, 0x95, 0x12, 0x00, 0x0F, 0x95, 0x12, 0x00,
			},
			TurnIndexBlocks{[]TurnIndexBlock{{From: 240751, Via: 1217808, To: 237761}, {From: 240751, Via: 1217808, To: 1217807}}, []byte{}},
		},
	}

	for _, c := range cases {
		var n TurnIndexBlocks
		writeLen, err := n.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p) != writeLen {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(n, c.n) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.n, n)
		}
	}
}
