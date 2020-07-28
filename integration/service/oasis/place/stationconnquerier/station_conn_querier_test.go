package stationconnquerier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/topograph"
)

func TestAppendIntoSortedSlice(t *testing.T) {
	cases := []struct {
		sortedArray      []*entity.RankedPlaceInfo
		itemToBeInserted *entity.RankedPlaceInfo
		expectedArray    []*entity.RankedPlaceInfo
	}{
		// case: insert into empty array
		{
			nil,
			&entity.RankedPlaceInfo{
				PlaceInfo: entity.PlaceInfo{
					ID:       3,
					Location: mockStation3Location,
				},
				Weight: &entity.Weight{
					Distance: 4622.08948420977,
					Duration: 208.2022290184581,
				},
			},
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
			},
		},

		// case: insert to the head of sorted array
		{
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &entity.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
			&entity.RankedPlaceInfo{
				PlaceInfo: entity.PlaceInfo{
					ID:       4,
					Location: mockStation4Location,
				},
				Weight: &entity.Weight{
					Distance: 222.0,
					Duration: 1.0,
				},
			},
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       4,
						Location: mockStation4Location,
					},
					Weight: &entity.Weight{
						Distance: 222.0,
						Duration: 1.0,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &entity.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
		},
		// case: insert into sorted array
		{
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &entity.Weight{
						Distance: 6310.598332634715,
						Duration: 284.2611861547169,
					},
				},
			},
			&entity.RankedPlaceInfo{
				PlaceInfo: entity.PlaceInfo{
					ID:       4,
					Location: mockStation4Location,
				},
				Weight: &entity.Weight{
					Distance: 4623.0,
					Duration: 1.0,
				},
			},
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       4,
						Location: mockStation4Location,
					},
					Weight: &entity.Weight{
						Distance: 4623.0,
						Duration: 1.0,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &entity.Weight{
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
		&topograph.MockConnectivityMap,
		mockOrigLocation,
		mockDestLocation,
		10,
		30,
	)

	// verify location
	locationCases := []struct {
		queryID        entity.PlaceID
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
			t.Errorf("Incorrect result for place.TopoQuerier.GetLocation, expect %+v but got %+v\n", c.expectLocation, actualLocation)
		}
	}

	// verify connectivity
	connectivityCases := []struct {
		placeID           entity.PlaceID
		expectQueryResult []*entity.RankedPlaceInfo
	}{
		{
			stationfindertype.OrigLocationID,
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 4622.08948420977,
						Duration: 208.2022290184581,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 4999.134247893073,
						Duration: 225.18622738257085,
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       1,
						Location: mockStation1Location,
					},
					Weight: &entity.Weight{
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
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       2,
						Location: mockStation2Location,
					},
					Weight: &entity.Weight{
						Distance: 1, // hard code value from mock ConnectivityMap
						Duration: 1, // hard code value from mock ConnectivityMap
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &entity.Weight{
						Distance: 4873.817197753869,
						Duration: 219.54131521413822,
					},
				},
			},
		},
		{
			3,
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &entity.Weight{
						Distance: 7083.8672907090095,
						Duration: 319.0931212031085,
					},
				},
			},
		},
		{
			2,
			[]*entity.RankedPlaceInfo{
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       3,
						Location: mockStation3Location,
					},
					Weight: &entity.Weight{
						Distance: 2, // hard code value from mock ConnectivityMap
						Duration: 2, // hard code value from mock ConnectivityMap
					},
				},
				{
					PlaceInfo: entity.PlaceInfo{
						ID:       stationfindertype.DestLocationID,
						Location: mockDestLocation,
					},
					Weight: &entity.Weight{
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
			t.Errorf("Incorrect result for place.TopoQuerier.NearByStationQuery, expect %#v but got %#v\n", c.expectQueryResult, actualQueryResult)
		}
	}
}

var mockPlaceInfo = []*entity.PlaceInfo{
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
func (finder *mockFinder) FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*entity.PlaceInfo {
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
