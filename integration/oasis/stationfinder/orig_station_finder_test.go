package stationfinder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

func createMockOrigStationFinder1() *origStationFinder {
	obj := &origStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse1,
			searchResp:        mockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func createMockOrigStationFinder2() *origStationFinder {
	obj := &origStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func createMockOrigStationFinder3() *origStationFinder {
	obj := &origStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func TestOrigStationFinderIterator(t *testing.T) {
	sf := createMockOrigStationFinder1()
	c := sf.iterateNearbyStations()
	var r []ChargeStationInfo
	go func() {
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo1) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo1, r)
		}
	}()
}
