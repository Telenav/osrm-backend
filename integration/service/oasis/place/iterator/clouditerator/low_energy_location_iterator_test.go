package clouditerator

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

func TestLowEnergyLocatioIterator1(t *testing.T) {
	sf := createMockLowEnergyLocationIterator1()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.IterateNearbyStations()
		var r []*iteratortype.ChargeStationInfo

		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo1) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo1, r)
		}
	}(&wg)
	wg.Wait()
}
func TestLowEnergyLocationIterator2(t *testing.T) {
	sf := createMockLowEnergyLocationIterator2()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.IterateNearbyStations()
		var r []*iteratortype.ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo2) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo2, r)
		}
	}(&wg)
	wg.Wait()
}
func TestLowEnergyLocationIterator3(t *testing.T) {
	sf := createMockLowEnergyLocationIterator3()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		c := sf.IterateNearbyStations()
		var r []*iteratortype.ChargeStationInfo
		for item := range c {
			r = append(r, item)
		}
		if !reflect.DeepEqual(r, mockChargeStationInfo3) {
			t.Errorf("expect %#v but got %#v", mockChargeStationInfo3, r)
		}
	}(&wg)
	wg.Wait()
}
