package builtinio

import (
	"reflect"
	"testing"
)

func TestBindWriterOnInt16Slice(t *testing.T) {
	cases := []struct {
		p []byte
		v []int16
	}{
		{
			[]byte{0x88, 0x03},
			[]int16{904},
		},
		{
			[]byte{
				0x88, 0x03,
				0xFF, // redundant byte
			},
			[]int16{904},
		},
		{
			[]byte{
				0x88, 0x03,
				0xFF, // redundant bytes
			},
			[]int16{904},
		},
		{
			[]byte{
				0x88, 0x03, 0xd2, 0x14,
			},
			[]int16{904, 5330},
		},
		{
			[]byte{0x88, 0x83},
			[]int16{-31864}, // negative
		},
	}

	for _, c := range cases {
		var u []int16
		w := BindWriterOnInt16Slice(&u)
		if w == nil {
			t.Errorf("BindWriterOnInt16Slice on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-l != len(c.p)%int16Bytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), l)
		}
		if !reflect.DeepEqual(u, c.v) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.v, u)
		}
	}

}

func TestBindWriterOnInt16SliceNil(t *testing.T) {
	var u *[]int16 // u == nil
	w := BindWriterOnInt16Slice(u)
	if w != nil {
		t.Errorf("BindWriterOnInt16Slice on nil, expect nil writer but got %+v", w)
	}
}
