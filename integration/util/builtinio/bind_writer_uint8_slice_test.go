package builtinio

import (
	"reflect"
	"testing"
)

func TestBindWriterOnUint8Slice(t *testing.T) {
	cases := []struct {
		p []byte
		v []uint8
	}{
		{
			[]byte{0x88},
			[]uint8{0x88},
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
			},
			[]uint8{0x88, 0x83, 0x00, 0x00},
		},
	}

	for _, c := range cases {
		var u []uint8
		w := BindWriterOnUint8Slice(&u)
		if w == nil {
			t.Errorf("BindWriterOnUint8Slice on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-l != 0 {
			t.Errorf("len(p) %d but write len %d", len(c.p), l)
		}
		if !reflect.DeepEqual(u, c.v) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.v, u)
		}
	}

}

func TestBindWriterOnUint8SliceNil(t *testing.T) {
	var u *[]uint8 // u == nil
	w := BindWriterOnUint8Slice(u)
	if w != nil {
		t.Errorf("BindWriterOnUint8Slice on nil, expect nil writer but got %+v", w)
	}
}
