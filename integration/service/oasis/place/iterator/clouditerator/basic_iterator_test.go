package clouditerator

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

var mockChargeStationInfo1 = []*iteratortype.ChargeStationInfo{
	{
		ID: "station1",
		Location: nav.Location{
			Lat: 32.333,
			Lon: 122.333,
		},
	},
	{
		ID: "station2",
		Location: nav.Location{
			Lat: -32.333,
			Lon: -122.333,
		},
	},
	{
		ID: "station3",
		Location: nav.Location{
			Lat: 32.333,
			Lon: -122.333,
		},
	},
	{
		ID: "station4",
		Location: nav.Location{
			Lat: -32.333,
			Lon: 122.333,
		},
	},
}

var mockChargeStationInfo2 = []*iteratortype.ChargeStationInfo{
	{
		ID: "station1",
		Location: nav.Location{
			Lat: 32.333,
			Lon: 122.333,
		},
	},
	{
		ID: "station2",
		Location: nav.Location{
			Lat: -32.333,
			Lon: -122.333,
		},
	},
	{
		ID: "station5",
		Location: nav.Location{
			Lat: -12.333,
			Lon: 122.333,
		},
	},
}

var mockChargeStationInfo3 = []*iteratortype.ChargeStationInfo{
	{
		ID: "station6",
		Location: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
	},
	{
		ID: "station7",
		Location: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
	},
}

func TestBasicIteratorCorrectness(t *testing.T) {
	cases := []struct {
		input  []*nearbychargestation.Result
		expect []*iteratortype.ChargeStationInfo
	}{
		{
			nearbychargestation.MockSearchResponse1.Results,
			mockChargeStationInfo1,
		},
	}

	for _, b := range cases {
		expect := b.expect
		bf := newBasicIterator(nil)
		bf.searchResp = &nearbychargestation.Response{}
		bf.searchResp.Results = b.input
		c := bf.IterateNearbyStations()

		var wg sync.WaitGroup
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()

			var r []*iteratortype.ChargeStationInfo
			for item := range c {
				r = append(r, item)
			}

			if !reflect.DeepEqual(r, expect) {
				t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
			}
		}(&wg)
		wg.Wait()
	}
}

func TestBasicIteratorAsync(t *testing.T) {
	cases := []struct {
		input     []*nearbychargestation.Result
		inputLock *sync.RWMutex
		expect    []*iteratortype.ChargeStationInfo
	}{
		{
			nearbychargestation.MockSearchResponse1.Results,
			&sync.RWMutex{},
			mockChargeStationInfo1,
		},
	}

	for _, b := range cases {
		expect := b.expect
		bf := newBasicIterator(nil)
		bf.searchResp = &nearbychargestation.Response{}
		bf.searchResp.Results = b.input
		num := 20
		var wg sync.WaitGroup
		for i := 0; i < num; i++ {
			go func(wg *sync.WaitGroup) {
				wg.Add(1)

				c := bf.IterateNearbyStations()
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					var r []*iteratortype.ChargeStationInfo
					for item := range c {
						r = append(r, item)
					}
					if !reflect.DeepEqual(r, expect) {
						t.Errorf("parse %v expect %v but got %v", b.input, b.expect, r)
					}
				}(wg)
			}(&wg)
		}
		wg.Wait()

	}
}
