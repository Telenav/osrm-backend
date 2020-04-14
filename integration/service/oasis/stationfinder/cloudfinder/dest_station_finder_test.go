package cloudfinder

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

func TestDestStationFinderIterator(t *testing.T) {
	sf := CreateMockDestStationFinder1()
	c := sf.IterateNearbyStations()
	var r []*stationfindertype.ChargeStationInfo

	for item := range c {
		r = append(r, item)
	}

	if !reflect.DeepEqual(r, mockChargeStationInfo1) {
		t.Errorf("expect %#v but got %#v", mockChargeStationInfo1, r)
	}
}
