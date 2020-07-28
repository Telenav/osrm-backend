package cloudfinder

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
)

func TestLowEnergyLocationStationFinderIterator1(t *testing.T) {
	sf := createMockLowEnergyLocationStationFinder1()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.IterateNearbyStations()
		var r []*stationfindertype.ChargeStationInfo

		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo1) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo1, r)
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

		c := sf.IterateNearbyStations()
		var r []*stationfindertype.ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo2) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo2, r)
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

		c := sf.IterateNearbyStations()
		var r []*stationfindertype.ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo3) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo3, r)
		}
	}(&wg)
	wg.Wait()
}
