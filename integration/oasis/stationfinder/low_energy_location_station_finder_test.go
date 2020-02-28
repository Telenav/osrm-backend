package stationfinder

import (
	"reflect"
	"sync"
	"testing"
)

func createMockLowEnergyLocationStationFinder1() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        mockSearchResponse1,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func createMockLowEnergyLocationStationFinder2() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        mockSearchResponse2,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func createMockLowEnergyLocationStationFinder3() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,
			searchResp:        mockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func TestLowEnergyLocationStationFinderIterator1(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder1()

	go func() {
		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo

		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo1) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo1, r)
		}
	}()
}
func TestLowEnergyLocationStationFinderIterator2(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder2()

	go func() {
		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo2) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo2, r)
		}
	}()
}
func TestLowEnergyLocationStationFinderIterator3(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder3()

	go func() {
		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo3) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo3, r)
		}
	}()
}
