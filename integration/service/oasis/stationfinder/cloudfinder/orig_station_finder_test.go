package cloudfinder

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

func TestOrigStationFinderIterator(t *testing.T) {
	sf := CreateMockOrigStationFinder1()
	c := sf.IterateNearbyStations()
	var r []*stationfindertype.ChargeStationInfo

	for item := range c {
		r = append(r, item)
	}

	if !reflect.DeepEqual(r, mockChargeStationInfo1) {
		t.Errorf("expect %#v but got %#v", mockChargeStationInfo1, r)
	}

}
