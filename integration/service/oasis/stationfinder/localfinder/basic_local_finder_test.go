package localfinder

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

func TestSingleIterator4BasicLocalFinder(t *testing.T) {
	localFinder := newBasicLocalFinder(nil)
	defer localFinder.Stop()

	mockFinder := spatialindexer.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPointIDs(spatialindexer.Location{}, 0, 0)

	iterC := localFinder.IterateNearbyStations()
	actual := make([]*stationfindertype.ChargeStationInfo, 0, len(stationfindertype.MockChargeStationInfo1))
	for item := range iterC {
		actual = append(actual, item)
	}

	if !reflect.DeepEqual(actual, stationfindertype.MockChargeStationInfo1) {
		t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", stationfindertype.MockChargeStationInfo1, actual)
	}

}

func TestMultipleIterator4BasicLocalFinder(t *testing.T) {
	localFinder := newBasicLocalFinder(nil)
	defer localFinder.Stop()

	mockFinder := spatialindexer.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPointIDs(spatialindexer.Location{}, 0, 0)

	for i := 0; i < 3; i++ {
		iterC := localFinder.IterateNearbyStations()
		actual := make([]*stationfindertype.ChargeStationInfo, 0, len(stationfindertype.MockChargeStationInfo1))
		for item := range iterC {
			actual = append(actual, item)
		}

		if !reflect.DeepEqual(actual, stationfindertype.MockChargeStationInfo1) {
			t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", stationfindertype.MockChargeStationInfo1, actual)
		}
	}

}
