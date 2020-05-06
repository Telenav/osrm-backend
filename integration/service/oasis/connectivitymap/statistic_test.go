package connectivitymap

import (
	"reflect"
	"testing"
)

var fakeDistanceLimit float64 = 123
var fakeStatisticResult1 = statistic{
	Count:                 5,
	ValidCount:            5,
	AverageNearByIDsCount: 2,
	MaxNearByIDsCount:     4,
	MinNearByIDsCount:     2,
	AverageMaxDistance:    15,
	MaxOfMaxDistance:      23,
	MinOfMaxDistance:      5,
	AverageMaxDuration:    24.8,
	MaxOfMaxDuration:      61,
	MinOfMaxDuration:      5,
	MaxRange:              123,
}

func TestStatisticBuild(t *testing.T) {
	cases := []struct {
		id2NearByIDsMap ID2NearByIDsMap
		distanceLimit   float64
		expect          *statistic
	}{
		{
			id2NearByIDsMap: fakeID2NearByIDsMap1,
			distanceLimit:   fakeDistanceLimit,
			expect:          &fakeStatisticResult1,
		},
		{
			id2NearByIDsMap: make(ID2NearByIDsMap),
			distanceLimit:   0.0,
			expect:          newStatistic(),
		},
	}

	for _, c := range cases {
		actual := newStatistic().build(c.id2NearByIDsMap, c.distanceLimit)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("Incorrect statistic build() result, expect \n%+v \nbut got \n%+v\n", c.expect, actual)
		}
	}
}
