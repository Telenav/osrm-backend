package stationconnquerier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

func TestAppendIntoSortedSlice(t *testing.T) {
	cases := []struct {
		sortedArray      []*common.RankedPlaceInfo
		itemToBeInserted *common.RankedPlaceInfo
		expectedArray    []*common.RankedPlaceInfo
	}{
		// case: insert into empty array
		{
			nil,
			&common.RankedPlaceInfo{
				PlaceInfo: common.PlaceInfo{
					ID:       3,
					Location: mockStation3Location,
				},
				Weight: &common.Weight{
					Distance: 4622.08948420977,
					Duration: 208.2022290184581,
				},
			},
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
			},
		},

		// case: insert to the head of sorted array
		{
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &common.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
			&common.RankedPlaceInfo{
				PlaceInfo: common.PlaceInfo{
					ID:       4,
					Location: mockStation4Location,
				},
				Weight: &common.Weight{
					Distance: 222.0,
					Duration: 1.0,
				},
			},
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       4,
						Location: mockStation4Location,
					},
					Weight: &common.Weight{
						Distance: 222.0,
						Duration: 1.0,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &common.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
		},
		// case: insert into sorted array
		{
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &common.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
			&common.RankedPlaceInfo{
				PlaceInfo: common.PlaceInfo{
					ID:       4,
					Location: mockStation4Location,
				},
				Weight: &common.Weight{
					Distance: 4623.0,
					Duration: 1.0,
				},
			},
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       4,
						Location: mockStation4Location,
					},
					Weight: &common.Weight{
						Distance: 4623.0,
						Duration: 1.0,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &common.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
		},
	}

	for _, c := range cases {
		actualResult := appendIntoSortedSlice(c.itemToBeInserted, c.sortedArray)
		if !reflect.DeepEqual(actualResult, c.expectedArray) {
			t.Errorf("Incorrect result expect %+v but got %+v\n", c.expectedArray, actualResult)
		}
	}

}

/*
Construct graph as follows

               station_1
            /      |      \
           /       |       \
  orig    ---   station_2   ---    dest
           \       |       /
            \      |      /
               station_3

Expects for connectivity:
orig: station_1, station_2, station_3
station_1: station_2, dest
station_2: station_3, dest
station_3: dest
dest: nil
*/
func TestStationConnQuerier(t *testing.T) {
	querier := New(
		&mockFinder{},
		ranker.CreateRanker(ranker.SimpleRanker, nil),
		&mockPlaceLocationQuerier{},
		&connectivitymap.MockConnectivityMap,
		mockOrigLocation,
		mockDestLocation,
		10,
		30,
	)

	// verify location
	locationCases := []struct {
		queryID        common.PlaceID
		expectLocation *nav.Location
	}{
		{
			stationfindertype.OrigLocationID,
			mockOrigLocation,
		},
		{
			stationfindertype.DestLocationID,
			mockDestLocation,
		},
		{
			1,
			mockStation1Location,
		},
		{
			2,
			mockStation2Location,
		},
		{
			3,
			mockStation3Location,
		},
		{
			stationfindertype.InvalidPlaceID,
			nil,
		},
	}

	for _, c := range locationCases {
		actualLocation := querier.GetLocation(c.queryID)
		if !reflect.DeepEqual(actualLocation, c.expectLocation) {
			t.Errorf("Incorrect result for connectivitymap.Querier.GetLocation, expect %+v but got %+v\n", c.expectLocation, actualLocation)
		}
	}

	// verify connectivity
	connectivityCases := []struct {
		placeID           common.PlaceID
		expectQueryResult []*common.RankedPlaceInfo
	}{
		{
			stationfindertype.OrigLocationID,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &common.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
		},
		{
			stationfindertype.DestLocationID,
			nil,
		},
		{
			1,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &common.Weight{
						Distance: 1, // hard code value from mock ConnectivityMap
						Duration: 1, // hard code value from mock ConnectivityMap
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &common.Weight{
						Distance: 4873.817197753869,
						Duration: 219.54131521413822,
					},
				},
			},
		},
		{
			3,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &common.Weight{
						Distance: 7083.8672907090095,
						Duration: 319.0931212031085,
					},
				},
			},
		},
		{
			2,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &common.Weight{
						Distance: 2, // hard code value from mock ConnectivityMap
						Duration: 2, // hard code value from mock ConnectivityMap
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &common.Weight{
						Distance: 7277.313067724465,
						Duration: 327.80689494254347,
					},
				},
			},
		},
	}

	for _, c := range connectivityCases {
		actualQueryResult := querier.NearByStationQuery(c.placeID)
		if !reflect.DeepEqual(actualQueryResult, c.expectQueryResult) {
			for _, r := range c.expectQueryResult {
				fmt.Printf("+++ %v, %v, %v, %v, %v\n", r.ID, r.Location.Lat, r.Location.Lon, r.Weight.Distance, r.Weight.Duration)
			}

			for _, r := range actualQueryResult {
				fmt.Printf("+++ %v, %v, %v, %v, %v\n", r.ID, r.Location.Lat, r.Location.Lon, r.Weight.Distance, r.Weight.Duration)
			}
			t.Errorf("Incorrect result for connectivitymap.Querier.NearByStationQuery, expect %#v but got %#v\n", c.expectQueryResult, actualQueryResult)
		}
	}
}

var mockPlaceInfo = []*common.PlaceInfo{
	{
		ID: 1,
		Location: &nav.Location{
			Lat: 37.355204,
			Lon: -121.953901,
		},
	},
	{
		ID: 2,
		Location: &nav.Location{
			Lat: 37.399331,
			Lon: -121.981193,
		},
	},
	{
		ID: 3,
		Location: &nav.Location{
			Lat: 37.401948,
			Lon: -121.977384,
		},
	},
}

type mockFinder struct {
}

// FindNearByPlaceIDs returns mock result
func (finder *mockFinder) FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*common.PlaceInfo {
	return mockPlaceInfo
}

type mockPlaceLocationQuerier struct {
}

var mockOrigLocation = &nav.Location{
	Lat: 37.407277,
	Lon: -121.925482,
}

var mockDestLocation = &nav.Location{
	Lat: 37.375024,
	Lon: -121.904706,
}

var mockStation1Location = &nav.Location{
	Lat: 37.355204,
	Lon: -121.953901,
}

var mockStation2Location = &nav.Location{
	Lat: 37.399331,
	Lon: -121.981193,
}

var mockStation3Location = &nav.Location{
	Lat: 37.401948,
	Lon: -121.977384,
}

var mockStation4Location = &nav.Location{
	Lat: 11.11,
	Lon: -22.22,
}

func (querier *mockPlaceLocationQuerier) GetLocation(placeID string) *nav.Location {
	switch placeID {
	case "1":
		return mockStation1Location
	case "2":
		return mockStation2Location
	case "3":
		return mockStation3Location
	}
	return nil
}
