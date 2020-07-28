package localiterator

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

func TestSingleIterator4BasicLocalIterator(t *testing.T) {
	localFinder := newBasicLocalIterator(nil)
	defer localFinder.Stop()

	mockFinder := mock.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPlaceIDs(nav.Location{}, 0, 0)

	iterC := localFinder.IterateNearbyStations()
	actual := make([]*iteratortype.ChargeStationInfo, 0, len(mock.MockChargeStationInfo1))
	for item := range iterC {
		actual = append(actual, item)
	}

	if !reflect.DeepEqual(actual, mock.MockChargeStationInfo1) {
		t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", mock.MockChargeStationInfo1, actual)
	}

}

func TestMultipleIterator4BasicLocalIterator(t *testing.T) {
	localFinder := newBasicLocalIterator(nil)
	defer localFinder.Stop()

	mockFinder := mock.MockFinder{}
	localFinder.placesInfo = mockFinder.FindNearByPlaceIDs(nav.Location{}, 0, 0)

	for i := 0; i < 3; i++ {
		iterC := localFinder.IterateNearbyStations()
		actual := make([]*iteratortype.ChargeStationInfo, 0, len(mock.MockChargeStationInfo1))
		for item := range iterC {
			actual = append(actual, item)
		}

		if !reflect.DeepEqual(actual, mock.MockChargeStationInfo1) {
			t.Errorf("Incorrect iterator result expect \n%#v\n but got \n%#v\n", mock.MockChargeStationInfo1, actual)
		}
	}

}
