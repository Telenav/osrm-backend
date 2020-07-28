package localfinder

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
)

func TestSingleIterator4BasicLocalFinder(t *testing.T) {
	localFinder := newBasicLocalFinder(nil)
	defer localFinder.Stop()

	mockFinder := mock.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPlaceIDs(nav.Location{}, 0, 0)

	iterC := localFinder.IterateNearbyStations()
	actual := make([]*stationfindertype.ChargeStationInfo, 0, len(mock.MockChargeStationInfo1))
	for item := range iterC {
		actual = append(actual, item)
	}

	if !reflect.DeepEqual(actual, mock.MockChargeStationInfo1) {
		t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", mock.MockChargeStationInfo1, actual)
	}

}

func TestMultipleIterator4BasicLocalFinder(t *testing.T) {
	localFinder := newBasicLocalFinder(nil)
	defer localFinder.Stop()

	mockFinder := mock.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPlaceIDs(nav.Location{}, 0, 0)

	for i := 0; i < 3; i++ {
		iterC := localFinder.IterateNearbyStations()
		actual := make([]*stationfindertype.ChargeStationInfo, 0, len(mock.MockChargeStationInfo1))
		for item := range iterC {
			actual = append(actual, item)
		}

		if !reflect.DeepEqual(actual, mock.MockChargeStationInfo1) {
			t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", mock.MockChargeStationInfo1, actual)
		}
	}

}
