package stationfinder

import (
	"reflect"
	"sync"
	"testing"
)

func createMockLowEnergyLocationStationFinder1() *destStationFinder {
	obj := &destStationFinder{
		osrmConnector:     nil,
		tnSearchConnector: nil,
		oasisReq:          nil,
		searchResp:        mockSearchResponse1,
		searchRespLock:    &sync.RWMutex{},
		bf:                &basicFinder{},
	}
	return obj
}

func createMockLowEnergyLocationStationFinder2() *origStationFinder {
	obj := &origStationFinder{
		osrmConnector:     nil,
		tnSearchConnector: nil,
		oasisReq:          nil,
		searchResp:        mockSearchResponse2,
		searchRespLock:    &sync.RWMutex{},
		bf:                &basicFinder{},
	}
	return obj
}

func createMockLowEnergyLocationStationFinder3() *origStationFinder {
	obj := &origStationFinder{
		osrmConnector:     nil,
		tnSearchConnector: nil,
		oasisReq:          nil,
		searchResp:        mockSearchResponse3,
		searchRespLock:    &sync.RWMutex{},
		bf:                &basicFinder{},
	}
	return obj
}

func TestLowEnergyLocationStationFinderIterator1(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder1()
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
func TestLowEnergyLocationStationFinderIterator2(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder2()
	c := sf.iterateNearbyStations()
	var r []ChargeStationInfo
	go func() {
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo2) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo1, r)
		}
	}()
}
func TestLowEnergyLocationStationFinderIterator3(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder3()
	c := sf.iterateNearbyStations()
	var r []ChargeStationInfo
	go func() {
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo3) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo1, r)
		}
	}()
}