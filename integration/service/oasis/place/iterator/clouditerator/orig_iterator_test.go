package clouditerator

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

func TestOrigIterator(t *testing.T) {
	sf := CreateMockOrigIterator1()
	c := sf.IterateNearbyStations()
	var r []*iteratortype.ChargeStationInfo

	for item := range c {
		r = append(r, item)
	}

	if !reflect.DeepEqual(r, mockChargeStationInfo1) {
		t.Errorf("expect %#v but got %#v", mockChargeStationInfo1, r)
	}

}
