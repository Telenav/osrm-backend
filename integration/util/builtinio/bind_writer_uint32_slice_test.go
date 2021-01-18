package builtinio

import (
	"reflect"
	"testing"
)

func TestBindWriterOnUint32Slice(t *testing.T) {
	cases := []struct {
		p []byte
		v []uint32
	}{
		{
			[]byte{0x88, 0x83, 0x00, 0x00},
			[]uint32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, // redundant byte
			},
			[]uint32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			[]uint32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00, 0xd2, 0x14, 0x11, 0x00,
			},
			[]uint32{33672, 1119442},
		},
	}

	for _, c := range cases {
		var u []uint32
		w := BindWriterOnUint32Slice(&u)
		if w == nil {
			t.Errorf("BindWriterOnUint32Slice on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-l != len(c.p)%uint32Bytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), l)
		}
		if !reflect.DeepEqual(u, c.v) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.v, u)
		}
	}

}

func TestBindWriterOnUint32SliceNil(t *testing.T) {
	var u *[]uint32 // u == nil
	w := BindWriterOnUint32Slice(u)
	if w != nil {
		t.Errorf("BindWriterOnUint32Slice on nil, expect nil writer but got %+v", w)
	}
}
