package builtinio

import (
	"testing"
)

func TestBindWriterOnUint32(t *testing.T) {
	cases := []struct {
		p []byte
		v uint32
	}{
		{[]byte{0x88, 0x83, 0x00, 0x00}, 33672},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, // redundant bytes
			},
			33672,
		},
		{
			[]byte{
				0x88, 0x83, 0x00, 0x00,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			33672,
		},
	}

	for _, c := range cases {
		var u uint32
		w := BindWriterOnUint32(&u)
		if w == nil {
			t.Errorf("BindWriterUint32 on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if l != uint32Bytes {
			t.Errorf("Write %v on %+v, expect write %d bytes but got %d bytes", c.p, u, uint32Bytes, l)
		}
		if u != c.v {
			t.Errorf("parse %v, expect to get %d but got %d", c.p, c.v, u)
		}
	}

}

func TestBindWriterOnUint32Nil(t *testing.T) {
	var u *uint32 // u == nil
	w := BindWriterOnUint32(u)
	if w != nil {
		t.Errorf("BindWriterUint32 on nil, expect nil writer but got %+v", w)
	}
}

func TestBindWriterOnUint32InsufficiantBytes(t *testing.T) {
	cases := []struct {
		p []byte
	}{
		{nil},
		{[]byte{}},
		{[]byte{0x88}},
		{[]byte{0x88, 0x83, 0x00}},
	}

	for _, c := range cases {
		var u uint32
		w := BindWriterOnUint32(&u)
		if w == nil {
			t.Errorf("BindWriterUint32 on %+v failed (got nil io.Writer)", &u)
		}
		_, err := w.Write(c.p)
		if err == nil {
			t.Errorf("Write %+v, expect error due to insufficiant bytes but got non error", c.p)
		}
	}

}
