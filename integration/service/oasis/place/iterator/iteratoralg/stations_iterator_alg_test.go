package iteratoralg

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/clouditerator"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/Telenav/osrm-backend/integration/util"
	"github.com/Telenav/osrm-backend/integration/util/osrmconnector"
)

var mockDict1 map[string]bool = map[string]bool{
	"station1": true,
	"station2": true,
	"station3": true,
	"station4": true,
}

func TestBuildChargeStationInfoDict1(t *testing.T) {
	sf := clouditerator.CreateMockOrigIterator1()
	m := buildChargeStationInfoDict(sf)
	if !reflect.DeepEqual(m, mockDict1) {
		t.Errorf("expect %v but got %v", mockDict1, m)
	}
}

var overlapChargeStationInfo1 = []*iteratortype.ChargeStationInfo{
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
}

func TestFindOverlapBetweenStations1(t *testing.T) {
	sf1 := clouditerator.CreateMockOrigIterator2()
	sf2 := clouditerator.CreateMockDestIterator1()
	r := FindOverlapBetweenStations(sf1, sf2)

	if !reflect.DeepEqual(r, overlapChargeStationInfo1) {
		t.Errorf("expect %v but got %v", overlapChargeStationInfo1, r)
	}
}

type fakeTableResponse struct {
}

func (ft *fakeTableResponse) Request4Table(r *table.Request) <-chan osrmconnector.TableResponse {
	c := make(chan osrmconnector.TableResponse)
	go func() {
		defer close(c)
		c <- osrmconnector.TableResponse{
			Resp: &mock.Mock4To2TableResponse1,
			Err:  nil,
		}
	}()
	return c
}

func TestCalcNeighborInfoPair(t *testing.T) {
	// from: station1, station2, station3, station4
	sf1 := clouditerator.CreateMockOrigIterator1()
	// to: station6, station7
	sf2 := clouditerator.CreateMockOrigIterator3()

	table := &fakeTableResponse{}
	r, err := CalcWeightBetweenChargeStationsPair(sf1, sf2, table)

	if err != nil {
		t.Errorf("expect no error but generate error of %v", err)
	}
	expect := []iteratortype.NeighborInfo{
		{
			FromID: "station1",
			FromLocation: nav.Location{
				Lat: 32.333,
				Lon: 122.333,
			},
			ToID: "station6",
			ToLocation: nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 2,
				Distance: 2,
			},
		},
		{
			FromID: "station1",
			FromLocation: nav.Location{
				Lat: 32.333,
				Lon: 122.333,
			},
			ToID: "station7",
			ToLocation: nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 3,
				Distance: 3,
			},
		},
		{
			FromID: "station2",
			FromLocation: nav.Location{
				Lat: -32.333,
				Lon: -122.333,
			},
			ToID: "station6",
			ToLocation: nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 4,
				Distance: 4,
			},
		},
		{
			FromID: "station2",
			FromLocation: nav.Location{
				Lat: -32.333,
				Lon: -122.333,
			},
			ToID: "station7",
			ToLocation: nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 5,
				Distance: 5,
			},
		},
		{
			FromID: "station3",
			FromLocation: nav.Location{
				Lat: 32.333,
				Lon: -122.333,
			},
			ToID: "station6",
			ToLocation: nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 6,
				Distance: 6,
			},
		},
		{
			FromID: "station3",
			FromLocation: nav.Location{
				Lat: 32.333,
				Lon: -122.333,
			},
			ToID: "station7",
			ToLocation: nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 7,
				Distance: 7,
			},
		},
		{
			FromID: "station4",
			FromLocation: nav.Location{
				Lat: -32.333,
				Lon: 122.333,
			},
			ToID: "station6",
			ToLocation: nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 8,
				Distance: 8,
			},
		},
		{
			FromID: "station4",
			FromLocation: nav.Location{
				Lat: -32.333,
				Lon: 122.333,
			},
			ToID: "station7",
			ToLocation: nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
			Weight: iteratortype.Weight{
				Duration: 9,
				Distance: 9,
			},
		},
	}
	if !reflect.DeepEqual(r, expect) {
		t.Errorf("expect %v but got %v", expect, r)
	}
}

// simulate locations array contains 4 points: orig -> location1 -> location2 -> dest
// (1.1, 1.1) -> (2.2, 2.2) -> (3.3, 3.3) -> (4.4, 4.4)
// location1 will find 4 nearby charge stations
// location2 will find 2 nearby charge stations
// search service will provide results based on upper information.
// Table service will provide result for: 1(orig) -> 4(charge stations around location 1),
// 4(charge stations around location 1) -> 2(charge stations around location 2),
// 2(charge stations around location 2) -> 1(dest)
func TestCalculateWeightBetweenNeighbors(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		if r.URL.EscapedPath() == "/entity/v4/search/json" {
			req, _ := nearbychargestation.ParseRequestURL(r.URL)
			if util.Float64Equal(req.Location.Lat, 2.2) && util.Float64Equal(req.Location.Lon, 2.2) {
				var searchResponseBytes4Location1, _ = json.Marshal(nearbychargestation.MockSearchResponse1)
				w.Write(searchResponseBytes4Location1)
			} else if util.Float64Equal(req.Location.Lat, 3.3) && util.Float64Equal(req.Location.Lon, 3.3) {
				var searchResponseBytes4Location2, _ = json.Marshal(nearbychargestation.MockSearchResponse3)
				w.Write(searchResponseBytes4Location2)
			}
			return
		}

		if strings.HasPrefix(r.URL.EscapedPath(), "/table/v1/driving/") {
			req, _ := table.ParseRequestURL(r.URL)
			s := len(req.Sources)
			d := len(req.Destinations)
			if s == 1 && d == 4 {
				var tableResponseBytesOrig2Location1, _ = json.Marshal(mock.Mock1To4TableResponse1)
				w.Write(tableResponseBytesOrig2Location1)
			} else if s == 4 && d == 2 {
				var tableResponseBytesLocation12Location2, _ = json.Marshal(mock.Mock4To2TableResponse1)
				w.Write(tableResponseBytesLocation12Location2)
			} else if s == 2 && d == 1 {
				var tableResponseBytesLocation2ToDest, _ = json.Marshal(mock.Mock2To1TableResponse1)
				w.Write(tableResponseBytesLocation2ToDest)
			}
			return
		}

	}))
	defer ts.Close()

	locations := []*nav.Location{
		{Lat: 1.1, Lon: 1.1},
		{Lat: 2.2, Lon: 2.2},
		{Lat: 3.3, Lon: 3.3},
		{Lat: 4.4, Lon: 4.4},
	}
	oc := osrmconnector.NewOSRMConnector(ts.URL)

	// create finder based on fake TNSearchService
	finder, err := iterator.CreateIteratorGenerator(iterator.CloudFinder, ts.URL, "apikey", "apisignature", nil)
	if err != nil {
		t.Errorf("Failed to create station finder during TestCalculateWeightBetweenNeighbors with error = %+v.\n", err)
	}
	c := CalculateWeightBetweenNeighbors(locations, oc, finder)

	expect_arr0 := mock.NeighborInfoArray0

	expect_arr1 := mock.NeighborInfoArray1

	expect_arr2 := mock.NeighborInfoArray2

	for arr := range c {
		switch len(arr.NeighborsInfo) {
		case 4:
			if !reflect.DeepEqual(arr.NeighborsInfo, expect_arr0) {
				t.Errorf("expect %v but got %v", expect_arr0, arr.NeighborsInfo)
			}
		case 8:
			if !reflect.DeepEqual(arr.NeighborsInfo, expect_arr1) {
				t.Errorf("expect %v but got %v", expect_arr1, arr.NeighborsInfo)
			}
		case 2:
			if !reflect.DeepEqual(arr.NeighborsInfo, expect_arr2) {
				t.Errorf("expect %v but got %v", expect_arr2, arr.NeighborsInfo)
			}
		}
	}
}
