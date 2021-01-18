package builtinio

import (
	"reflect"
	"testing"
)

func TestBindWriterOnInt32Slice(t *testing.T) {
	cases := []struct {
		p []byte
		v []int32
	}{
		{
			[]byte{0x88, 0x83, 0x00, 0x00},
			[]int32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, // redundant byte
			},
			[]int32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			[]int32{33672},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00, 0xd2, 0x14, 0x11, 0x00,
			},
			[]int32{33672, 1119442},
		},
		{
			[]byte{0xE8, 0x05, 0x00, 0x80},
			[]int32{-2147482136}, // negative
		},
	}

	for _, c := range cases {
		var u []int32
		w := BindWriterOnInt32Slice(&u)
		if w == nil {
			t.Errorf("BindWriterOnInt32Slice on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-l != len(c.p)%int32Bytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), l)
		}
		if !reflect.DeepEqual(u, c.v) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.v, u)
		}
	}

}

func TestBindWriterOnInt32SliceNil(t *testing.T) {
	var u *[]int32 // u == nil
	w := BindWriterOnInt32Slice(u)
	if w != nil {
		t.Errorf("BindWriterOnUint32Slice on nil, expect nil writer but got %+v", w)
	}
}
