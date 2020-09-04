package osrmtype

import (
	"reflect"
	"testing"
)

func TestWriteCoordinates(t *testing.T) {
	cases := []struct {
		p []byte
		Coordinates
	}{
		{
			[]byte{
				0x19, 0xc4, 0xd8, 0xf8, 0x8c, 0xdb, 0x59, 0x02,
			},
			Coordinates{
				Coordinate{FixedLon: -120011751, FixedLat: 39443340},
			},
		},
		{
			[]byte{
				0x19, 0xc4, 0xd8, 0xf8, 0x8c, 0xdb, 0x59, 0x02,
				0xFF, // redundant byte
			},
			Coordinates{
				Coordinate{FixedLon: -120011751, FixedLat: 39443340},
			},
		},
		{
			[]byte{
				0x19, 0xc4, 0xd8, 0xf8, 0x8c, 0xdb, 0x59, 0x02,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // redundant bytes
			},
			Coordinates{
				Coordinate{FixedLon: -120011751, FixedLat: 39443340},
			},
		},
		{
			[]byte{
				0x19, 0xc4, 0xd8, 0xf8, 0x8c, 0xdb, 0x59, 0x02,
				0x66, 0x03, 0x24, 0xf9, 0x6b, 0xba, 0x27, 0x02,
			},
			Coordinates{
				Coordinate{FixedLon: -120011751, FixedLat: 39443340},
				Coordinate{FixedLon: -115080346, FixedLat: 36158059},
			},
		},
	}

	for _, c := range cases {
		coordinates := Coordinates{}
		writeLen, err := coordinates.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if len(c.p)-writeLen != len(c.p)%CoordinateBytes {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		if !reflect.DeepEqual(coordinates, c.Coordinates) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.Coordinates, coordinates)
		}
	}
}
