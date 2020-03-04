package stationfinder

import (
	"reflect"
	"sync"
	"testing"
  
	"github.com/Telenav/osrm-backend/integration/pkg/api/search/nearbychargestation"
)

func createMockLowEnergyLocationStationFinder1() *lowEnergyLocationStationFinder {
	obj := &lowEnergyLocationStationFinder{
		location: nil,
		bf: &basicFinder{
			tnSearchConnector: nil,

			searchResp:        nearbychargestation.MockSearchResponse1,
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
			searchResp:        nearbychargestation.MockSearchResponse2,
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
			searchResp:        nearbychargestation.MockSearchResponse3,
			searchRespLock:    &sync.RWMutex{},
		},
	}
	return obj
}

func TestLowEnergyLocationStationFinderIterator1(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder1()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo

		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo1) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo1, r)
		}
	}(&wg)
	wg.Wait()
}
func TestLowEnergyLocationStationFinderIterator2(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder2()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo2) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo2, r)
		}
	}(&wg)
	wg.Wait()
}
func TestLowEnergyLocationStationFinderIterator3(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder3()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.iterateNearbyStations()
		var r []ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo3) {
			t.Errorf("expect %v but got %v", mockChargeStationInfo3, r)
		}
	}(&wg)
	wg.Wait()
}
