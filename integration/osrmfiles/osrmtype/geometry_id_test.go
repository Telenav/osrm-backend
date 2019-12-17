package osrmtype

import (
	"math"
	"reflect"
	"testing"
)

func TestPraseGeometryID(t *testing.T) {
	cases := []struct {
		p []byte
		GeometryID
	}{
		// Use MaxUint >> 1 as invalid ID of GeometryID
		// See https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/util/typedefs.hpp#L159
		{[]byte{0xff, 0xff, 0xff, 0x7f}, GeometryID{NodeID(math.MaxUint32 >> 1), false}},
	}

	for _, c := range cases {
		id := GeometryID{}
		if err := id.tryParse(c.p); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(id, c.GeometryID) {
			t.Errorf("parse %v expect GeometryID %v but got %v", c.p, c.GeometryID, id)
		}
	}
}

func TestPraseGeometryIDFail(t *testing.T) {
	cases := []struct {
		p []byte
	}{
		{nil},
		{[]byte{}},
		{[]byte{0xff, 0xff, 0xff}},
	}

	for _, c := range cases {
		id := GeometryID{}
		if err := id.tryParse(c.p); err == nil {
			t.Errorf("parse %v expect fail but succeed, got %v", c.p, id)
		}
	}
}
