package builtinio

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/util"
)

func TestBindWriterOnFloat3232Slice(t *testing.T) {
	cases := []struct {
		p []byte
		v []float32
	}{
		{
			[]byte{0x6D, 0xA7, 0x4F, 0x44},
			[]float32{830.616028},
		},
		{
			[]byte{
				0x6D, 0xA7, 0x4F, 0x44,
				0xFF, // redundant byte
			},
			[]float32{830.616028},
		},
		{
			[]byte{
				0x6D, 0xA7, 0x4F, 0x44,
				0xFF, 0xFF, 0xFF, // redundant bytes
			},
			[]float32{830.616028},
		},
		{
			[]byte{
				0x6D, 0xA7, 0x4F, 0x44, 0xB3, 0xA2, 0x8C, 0x42,
			},
			[]float32{830.616028, 70.3177719},
		},
		{
			[]byte{0x68, 0xCD, 0x8F, 0xBF},
			[]float32{-1.123456}, // negative
		},
	}

	for _, c := range cases {
		var u []float32
		w := BindWriterOnFloat32Slice(&u)
		if w == nil {
			t.Errorf("BindWriterOnFloat32Slice on %+v failed (got nil io.Writer)", &u)
		}
		l, err := w.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-l != len(c.p)%float32Bytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), l)
		}
		for i := range c.v {
			if !util.Float32Equal(c.v[i], u[i]) {
				t.Errorf("parse %v expect %v but got %v", c.p, c.v, u)
			}
		}
	}

}

func TestBindWriterOnFloat32SliceNil(t *testing.T) {
	var u *[]float32 // u == nil
	w := BindWriterOnFloat32Slice(u)
	if w != nil {
		t.Errorf("BindWriterOnFloat32Slice on nil, expect nil writer but got %+v", w)
	}
}
