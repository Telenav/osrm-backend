package packed

import (
	"reflect"
	"testing"
)

func TestWriteUint64Vector(t *testing.T) {

	packed63BitsBlock := []byte{ // 63 bits, totally 63*8 bytes per block
		0x85, 0xea, 0x98, 0x01, 0x00, 0x00, 0x00, 0x00, 0x43, 0x75, 0xcc, 0x00, 0x00, 0x00, 0x00, 0xc0,
		0xa1, 0x3a, 0x66, 0x00, 0x00, 0x00, 0x00, 0x00, 0x51, 0x1d, 0x33, 0x00, 0x00, 0x00, 0x00, 0x90,
		0xa8, 0x8e, 0x19, 0x00, 0x00, 0x00, 0x00, 0x50, 0x54, 0xc7, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x64,
		0xde, 0x64, 0x06, 0x00, 0x00, 0x00, 0x00, 0x34, 0x6f, 0x32, 0x03, 0x00, 0x00, 0x00, 0x00, 0x9b,
		0x37, 0x99, 0x01, 0x00, 0x00, 0x00, 0x00, 0xce, 0x9b, 0xcc, 0x00, 0x00, 0x00, 0x00, 0xc0, 0xe7,
		0x4d, 0x66, 0x00, 0x00, 0x00, 0x00, 0x40, 0xe7, 0x2d, 0x3c, 0x00, 0x00, 0x00, 0x00, 0xb0, 0xf3,
		0x16, 0x1e, 0x00, 0x00, 0x00, 0x00, 0x88, 0x7b, 0x0b, 0x0f, 0x00, 0x00, 0x00, 0x00, 0xc8, 0xbd,
		0x85, 0x07, 0x00, 0x00, 0x00, 0x00, 0xf4, 0xe1, 0xc2, 0x03, 0x00, 0x00, 0x00, 0x00, 0xfb, 0x70,
		0xe1, 0x01, 0x00, 0x00, 0x00, 0x80, 0x52, 0xb8, 0x09, 0x01, 0x00, 0x00, 0x00, 0x40, 0x75, 0xdc,
		0x84, 0x00, 0x00, 0x00, 0x00, 0xe0, 0xdd, 0x6f, 0x42, 0x00, 0x00, 0x00, 0x00, 0x10, 0x6a, 0x38,
		0x21, 0x00, 0x00, 0x00, 0x00, 0x38, 0x3f, 0x9c, 0x10, 0x00, 0x00, 0x00, 0x00, 0x60, 0x24, 0x4e,
		0x08, 0x00, 0x00, 0x00, 0x00, 0x5a, 0x23, 0x27, 0x04, 0x00, 0x00, 0x00, 0x00, 0xec, 0x97, 0x13,
		0x02, 0x00, 0x00, 0x00, 0x80, 0xf4, 0xce, 0x09, 0x01, 0x00, 0x00, 0x00, 0x80, 0x90, 0xe8, 0x84,
		0x00, 0x00, 0x00, 0x00, 0xe0, 0x56, 0x74, 0x42, 0x00, 0x00, 0x00, 0x00, 0x30, 0xb5, 0x3a, 0x21,
		0x00, 0x00, 0x00, 0x00, 0x68, 0x9c, 0x9d, 0x10, 0x00, 0x00, 0x00, 0x00, 0xec, 0xeb, 0x4e, 0x08,
		0x00, 0x00, 0x00, 0x00, 0x38, 0x8b, 0x27, 0x04, 0x00, 0x00, 0x00, 0x00, 0xd5, 0xcd, 0x13, 0x02,
		0x00, 0x00, 0x00, 0x80, 0x50, 0xeb, 0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0xf7, 0xf6, 0x84, 0x00,
		0x00, 0x00, 0x00, 0xc0, 0xfe, 0x7c, 0x42, 0x00, 0x00, 0x00, 0x00, 0x30, 0xfd, 0x3e, 0x21, 0x00,
		0x00, 0x00, 0x00, 0xd8, 0xbd, 0x9f, 0x10, 0x00, 0x00, 0x00, 0x00, 0x4c, 0xf3, 0x4f, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x4c, 0x10, 0x28, 0x04, 0x00, 0x00, 0x00, 0x00, 0xa6, 0x0e, 0x14, 0x02, 0x00,
		0x00, 0x00, 0x80, 0x8a, 0x07, 0x0a, 0x01, 0x00, 0x00, 0x00, 0xc0, 0xb3, 0x04, 0x85, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x55, 0x83, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x42, 0x21, 0x00, 0x00,
		0x00, 0x00, 0x50, 0x17, 0xa1, 0x10, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x9b, 0x50, 0x08, 0x00, 0x00,
		0x00, 0x00, 0x6a, 0x52, 0x28, 0x04, 0x00, 0x00, 0x00, 0x00, 0x1d, 0x2a, 0x14, 0x02, 0x00, 0x00,
		0x00, 0x80, 0x3c, 0x17, 0x0a, 0x01, 0x00, 0x00, 0x00, 0x80, 0x30, 0x0c, 0x85, 0x00, 0x00, 0x00,
		0x00, 0x20, 0x53, 0x86, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x43, 0x21, 0x00, 0x00, 0x00,
		0x00, 0x38, 0x05, 0xa2, 0x10, 0x00, 0x00, 0x00, 0x00, 0x84, 0x18, 0x51, 0x08, 0x00, 0x00, 0x00,
		0x00, 0xea, 0x92, 0x28, 0x04, 0x00, 0x00, 0x00, 0x00, 0x6b, 0x4d, 0x14, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x29, 0x2d, 0x0a, 0x01, 0x00, 0x00, 0x00, 0x40, 0xb8, 0x17, 0x85, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x87, 0x8c, 0x42, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x82, 0x46, 0x21, 0x00, 0x00, 0x00, 0x00,
		0x70, 0x5e, 0xa3, 0x10, 0x00, 0x00, 0x00, 0x00, 0xc8, 0xd4, 0x51, 0x08, 0x00, 0x00, 0x00, 0x00,
		0x42, 0xf9, 0x28, 0x04, 0x00, 0x00, 0x00, 0x00,
	}
	packed63BitsVec := NewUint64Vector(63)
	packed63BitsVec.NumOfElements = 64
	packed63BitsVec.PackedMeta = 63
	packed63BitsVec.Values = []uint64{
		26798725, 26798726, 26798727, 26798728, 26798729, 26798730, 26818457, 26818458, 26818459, 26818460,
		26818463, 31551290, 31551291, 31551345, 31551346, 31551738, 31551739, 34828453, 34828757, 34832111,
		34834081, 34834407, 34834712, 34836909, 34838508, 34840041, 34841154, 34841271, 34843475, 34845581,
		34847483, 34850204, 34852309, 34854561, 34855900, 34858998, 34861011, 34863035, 34864339, 34867238,
		34868902, 34869013, 34869967, 34871976, 34873616, 34874090, 34875075, 34875701, 34875933, 34877049,
		34877634, 34878105, 34878432, 34881703, 34883105, 34883957, 34884971, 34888274, 34889441, 34890808,
		34891819, 34892750, 34895154, 34897057,
	}

	cases := []struct {
		p []byte
		Uint64Vector
	}{
		{
			packed63BitsBlock, packed63BitsVec,
		},
	}

	for _, c := range cases {
		u := NewUint64Vector(c.Uint64Vector.bits)
		writeLen, err := u.Write(c.p)
		if err != nil {
			t.Error(err)
		}
		if writeLen != len(c.p) {
			t.Errorf("len(p) %d but write len %d", len(c.p), writeLen)
		}
		u.NumOfElements = c.Uint64Vector.NumOfElements
		u.PackedMeta = c.Uint64Vector.PackedMeta
		if err := u.Prune(); err != nil {
			t.Error(err)
		}
		if err := u.Validate(); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(u.Values, c.Uint64Vector.Values) {
			t.Errorf("parse %v expect %v but got %v", c.p, c.Uint64Vector, u)
		}
	}

}
