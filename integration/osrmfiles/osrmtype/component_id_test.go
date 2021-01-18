package osrmtype

import (
	"math"
	"reflect"
	"testing"
)

func TestPraseComponentID(t *testing.T) {
	cases := []struct {
		p []byte
		ComponentID
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff}, ComponentID{math.MaxUint32 >> 1, true}},
		{[]byte{0xff, 0xff, 0xff, 0x7f}, ComponentID{math.MaxUint32 >> 1, false}},
	}

	for _, c := range cases {
		id := ComponentID{}
		if err := id.tryParse(c.p); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(id, c.ComponentID) {
			t.Errorf("parse %v expect ComponentID %v but got %v", c.p, c.ComponentID, id)
		}
	}
}

func TestPraseComponentIDFail(t *testing.T) {
	cases := []struct {
		p []byte
	}{
		{nil},
		{[]byte{}},
		{[]byte{0xff, 0xff, 0xff}},
	}

	for _, c := range cases {
		id := ComponentID{}
		if err := id.tryParse(c.p); err == nil {
			t.Errorf("parse %v expect fail but succeed, got %v", c.p, id)
		}
	}
}
