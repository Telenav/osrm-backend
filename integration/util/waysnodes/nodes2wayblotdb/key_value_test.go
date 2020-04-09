package nodes2wayblotdb

import (
	"reflect"
	"testing"
)

func TestKey(t *testing.T) {
	cases := []struct {
		fromNodeID int64
		toNodeID   int64
		key        []byte
	}{
		{0, 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{24418325, 84760891102, []byte{0x15, 0x98, 0x74, 0x01, 0x00, 0x00, 0x00, 0x00, 0xDE, 0x8E, 0x24, 0xBC, 0x13, 0x00, 0x00, 0x00}},
	}

	for _, c := range cases {
		k := key(c.fromNodeID, c.toNodeID)
		if len(k) == 0 || !reflect.DeepEqual(k, c.key) {
			t.Errorf("generate key from %d,%d, expect to have %v, but got %v", c.fromNodeID, c.toNodeID, c.key, k)
		}

		from, to := parseKey(c.key)
		if from != c.fromNodeID || to != c.toNodeID {
			t.Errorf("parse key from %v, expect %d,%d but got %d,%d", c.key, c.fromNodeID, c.toNodeID, from, to)
		}
	}

}

func TestValue(t *testing.T) {
	cases := []struct {
		wayID int64
		value []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{84760891102, []byte{0xDE, 0x8E, 0x24, 0xBC, 0x13, 0x00, 0x00, 0x00}},
	}

	for _, c := range cases {
		v := value(c.wayID)
		if len(v) == 0 || !reflect.DeepEqual(v, c.value) {
			t.Errorf("generate value from %d, expect to have %v, but got %v", c.wayID, c.value, v)
		}

		wayID := parseValue(c.value)
		if wayID != c.wayID {
			t.Errorf("parse value from %v, expect %d but got %d", c.value, c.wayID, wayID)
		}
	}

}
