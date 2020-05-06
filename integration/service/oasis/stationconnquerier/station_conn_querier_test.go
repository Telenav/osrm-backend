package stationconnquerier

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

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
		queryStr       string
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
			"1",
			mockStation1Location,
		},
		{
			"2",
			mockStation2Location,
		},
		{
			"3",
			mockStation3Location,
		},
		{
			"incorrect_station_id",
			nil,
		},
	}

	for _, c := range locationCases {
		actualLocation := querier.GetLocation(c.queryStr)
		if !reflect.DeepEqual(actualLocation, c.expectLocation) {
			t.Errorf("Incorrect result for connectivitymap.Querier.GetLocation, expect %+v but got %+v\n", c.expectLocation, actualLocation)
		}
	}

	// verify connectivity
	connectivityCases := []struct {
		stationID         string
		expectQueryResult []*connectivitymap.QueryResult
	}{
		{
			stationfindertype.OrigLocationID,
			[]*connectivitymap.QueryResult{
				{
					StationID:       "3",
					StationLocation: mockStation3Location,
					Distance:        4622.08948420977,
					Duration:        4622.08948420977,
				},
				{
					StationID:       "2",
					StationLocation: mockStation2Location,
					Distance:        4999.134247893073,
					Duration:        4999.134247893073,
				},
				{
					StationID:       "1",
					StationLocation: mockStation1Location,
					Distance:        6310.598332634715,
					Duration:        6310.598332634715,
				},
			},
		},
		{
			stationfindertype.DestLocationID,
			nil,
		},
		{
			"1",
			[]*connectivitymap.QueryResult{
				{
					StationID:       "2",
					StationLocation: mockStation2Location,
					Distance:        1, // hard code value from mock ConnectivityMap
					Duration:        1, // hard code value from mock ConnectivityMap
				},
				{
					StationID:       stationfindertype.DestLocationID,
					StationLocation: mockDestLocation,
					Distance:        4873.817197753869,
					Duration:        219.54131521413822,
				},
			},
		},
		{
			"2",
			[]*connectivitymap.QueryResult{
				{
					StationID:       "3",
					StationLocation: mockStation3Location,
					Distance:        2, // hard code value from mock ConnectivityMap
					Duration:        2, // hard code value from mock ConnectivityMap
				},
				{
					StationID:       stationfindertype.DestLocationID,
					StationLocation: mockDestLocation,
					Distance:        7277.313067724465,
					Duration:        327.80689494254347,
				},
			},
		},
		{
			"3",
			[]*connectivitymap.QueryResult{
				{
					StationID:       stationfindertype.DestLocationID,
					StationLocation: mockDestLocation,
					Distance:        7083.8672907090095,
					Duration:        319.0931212031085,
				},
			},
		},
	}

	for _, c := range connectivityCases {
		actualQueryResult := querier.NearByStationQuery(c.stationID)
		if !reflect.DeepEqual(actualQueryResult, c.expectQueryResult) {
			t.Errorf("Incorrect result for connectivitymap.Querier.NearByStationQuery, expect %+v but got %+v\n", c.expectQueryResult, actualQueryResult)
		}
	}
}

var mockPlaceInfo = []*spatialindexer.PointInfo{
	{
		ID: 1,
		Location: spatialindexer.Location{
			Lat: 37.355204,
			Lon: -121.953901,
		},
	},
	{
		ID: 2,
		Location: spatialindexer.Location{
			Lat: 37.399331,
			Lon: -121.981193,
		},
	},
	{
		ID: 3,
		Location: spatialindexer.Location{
			Lat: 37.401948,
			Lon: -121.977384,
		},
	},
}

type mockFinder struct {
}

// FindNearByPointIDs returns mock result
func (finder *mockFinder) FindNearByPointIDs(center spatialindexer.Location, radius float64, limitCount int) []*spatialindexer.PointInfo {
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
