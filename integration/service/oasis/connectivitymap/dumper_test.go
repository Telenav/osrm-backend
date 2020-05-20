package connectivitymap

import (
	"os"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/golang/glog"
)

func TestDumpGivenObjectThenLoadAndThenCompareWithOriginalObject(t *testing.T) {
	cases := []ConnectivityMap{
		ConnectivityMap{
			id2nearByIDs: fakeID2NearByIDsMap1,
			maxRange:     fakeDistanceLimit,
			statistic:    &fakeStatisticResult1,
		},
	}

	// check whether curent folder is writeable
	path, _ := os.Getwd()
	_, err := os.OpenFile(path+"/tmp", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	if err := os.Remove(path + "/tmp"); err != nil {
		glog.Errorf("During dumper test, remove path %s failed.\n", path+"/tmp")
		return
	}

	if err := removeAllDumpFiles(path); err != nil {
		t.Errorf("During running removeAllDumpFiles met error %v", err)
	}
	for _, c := range cases {
		if err := serializeConnectivityMap(&c, path); err != nil {
			t.Errorf("During running serializeConnectivityMap for case %v, met error %v", c, err)
		}

		actual := New(0.0)
		if err := deSerializeConnectivityMap(actual, path); err != nil {
			t.Errorf("During running deSerializeConnectivityMap for case %v, met error %v", c, err)
		}

		if !reflect.DeepEqual(actual.id2nearByIDs, c.id2nearByIDs) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.id2nearByIDs, actual.id2nearByIDs)
		}

		if !reflect.DeepEqual(actual.maxRange, c.maxRange) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.maxRange, actual.maxRange)
		}

		if !reflect.DeepEqual(actual.statistic, c.statistic) {
			t.Errorf("Expect result \n%+v but got \n%+v\n", c.statistic, actual.statistic)
		}

		if err := removeAllDumpFiles(path); err != nil {
			t.Errorf("During running removeAllDumpFiles met error %v", err)
		}
	}
}

/*
fakeID2NearByIDsMap1 represents following station graph:

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
var fakeID2NearByIDsMap1 = ID2NearByIDsMap{
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
