package osrmtype

import (
	"reflect"
	"testing"
)

func TestPraseRoadClassification(t *testing.T) {
	cases := []struct {
		p []byte
		RoadClassification
	}{
		{
			[]byte{0x01, 0x02},
			RoadClassification{
				MotowayClass:      true,
				LinkClass:         false,
				MaybeIgnored:      false,
				RoadPriorityClass: RoadPriorityClassMotoway,
				NumberOfLanes:     2,
			}},
	}

	for _, c := range cases {
		r := RoadClassification{}
		if err := r.tryParse(c.p); err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(r, c.RoadClassification) {
			t.Errorf("parse %v expect RoadClassification %v but got %v", c.p, c.RoadClassification, r)
		}
	}
}

func TestPraseRoadClassificationFail(t *testing.T) {
	cases := []struct {
		p []byte
	}{
		{nil},
		{[]byte{}},
		{[]byte{0xff}},
	}

	for _, c := range cases {
		r := RoadClassification{}
		if err := r.tryParse(c.p); err == nil {
			t.Errorf("parse %v expect fail but succeed, got %v", c.p, r)
		}
	}
}
