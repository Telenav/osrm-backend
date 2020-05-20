package connectivitymap

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
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
			id2NearByIDsMap: fakeID2NearByIDsMap2,
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

/*
fakeID2NearByIDsMap2 represents following station graph:

7             5
|    \   /    |
|      1      |
|    /   \    |
2             8

1 -> 2(Distance: 3, Duration: 3),
	 5(Distance: 4, Duration: 4),
	 7(Distance: 6, Duration: 61),
	 8(Distance:12, Duration: 12),
2 -> 1(Distance: 3, Duration: 3),
	 7(Distance: 23, Duration: 23),
5 -> 1(Distance: 4, Duration: 4),
	 8(Distance: 5, Duration: 5),
7 -> 1(Distance: 6, Duration: 6),
	 2(Distance: 23, Duration: 23),
8 -> 5(Distance: 5, Duration: 5),
     1(Distance: 12, Duration: 12),
*/
var fakeID2NearByIDsMap2 = ID2NearByIDsMap{
	1: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 2,
			},
			Weight: &common.Weight{
				Distance: 3,
				Duration: 3,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 5,
			},
			Weight: &common.Weight{
				Distance: 4,
				Duration: 4,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 7,
			},
			Weight: &common.Weight{
				Distance: 6,
				Duration: 61,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 8,
			},
			Weight: &common.Weight{
				Distance: 12,
				Duration: 12,
			},
		},
	},

	2: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 1,
			},
			Weight: &common.Weight{
				Distance: 3,
				Duration: 3,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 7,
			},
			Weight: &common.Weight{
				Distance: 23,
				Duration: 23,
			},
		},
	},

	5: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 1,
			},
			Weight: &common.Weight{
				Distance: 4,
				Duration: 4,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 8,
			},
			Weight: &common.Weight{
				Distance: 5,
				Duration: 5,
			},
		},
	},

	7: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 1,
			},
			Weight: &common.Weight{
				Distance: 6,
				Duration: 6,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 2,
			},
			Weight: &common.Weight{
				Distance: 23,
				Duration: 23,
			},
		},
	},

	8: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 5,
			},
			Weight: &common.Weight{
				Distance: 5,
				Duration: 5,
			},
		},
		{
			PlaceInfo: common.PlaceInfo{
				ID: 1,
			},
			Weight: &common.Weight{
				Distance: 12,
				Duration: 12,
			},
		},
	},
}
