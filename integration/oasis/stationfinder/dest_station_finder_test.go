package stationfinder

import (
	"reflect"
	"sync"
	"testing"
)

func createMockDestStationFinder1() *destStationFinder {
	obj := &destStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        mockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func createMockDestStationFinder2() *destStationFinder {
	obj := &destStationFinder{
		oasisReq: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        mockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func TestDestStationFinderIterator(t *testing.T) {
	sf := createMockDestStationFinder1()
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
