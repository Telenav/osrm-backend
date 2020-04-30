package stationfinderalg

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

// orig_location -> station1, station2, station3, station4
// station1 -> station6, station7
// station2 -> station6, station7
// station3 -> station6, station7
// station4 -> station6, station7
// station6 -> destination
// station7 -> destination
func TestQuerierBasedOnWeightBetweenNeighborsChan(t *testing.T) {
	cases := []struct {
		stationID         string
		expectQueryResult []*connectivitymap.QueryResult
		expectLocation    *nav.Location
	}{
		{
			"orig_location",
			[]*connectivitymap.QueryResult{
				{
					StationID: "station1",
					StationLocation: &nav.Location{
						Lat: 32.333,
						Lon: 122.333,
					},
					Distance: 22.2,
					Duration: 22.2,
				},
				{
					StationID: "station2",
					StationLocation: &nav.Location{
						Lat: -32.333,
						Lon: -122.333,
					},
					Distance: 11.1,
					Duration: 11.1,
				},
				{
					StationID: "station3",
					StationLocation: &nav.Location{
						Lat: 32.333,
						Lon: -122.333,
					},
					Distance: 33.3,
					Duration: 33.3,
				},
				{
					StationID: "station4",
					StationLocation: &nav.Location{
						Lat: -32.333,
						Lon: 122.333,
					},
					Distance: 44.4,
					Duration: 44.4,
				},
			},
			&nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			"station1",
			[]*connectivitymap.QueryResult{
				{
					StationID: "station6",
					StationLocation: &nav.Location{
						Lat: 30.333,
						Lon: 122.333,
					},
					Distance: 2,
					Duration: 2,
				},
				{
					StationID: "station7",
					StationLocation: &nav.Location{
						Lat: -10.333,
						Lon: 122.333,
					},
					Distance: 3,
					Duration: 3,
				},
			},
			&nav.Location{
				Lat: 32.333,
				Lon: 122.333,
			},
		},
		{
			"station2",
			[]*connectivitymap.QueryResult{
				{
					StationID: "station6",
					StationLocation: &nav.Location{
						Lat: 30.333,
						Lon: 122.333,
					},
					Distance: 4,
					Duration: 4,
				},
				{
					StationID: "station7",
					StationLocation: &nav.Location{
						Lat: -10.333,
						Lon: 122.333,
					},
					Distance: 5,
					Duration: 5,
				},
			},
			&nav.Location{
				Lat: -32.333,
				Lon: -122.333,
			},
		},
		{
			"station3",
			[]*connectivitymap.QueryResult{
				{
					StationID: "station6",
					StationLocation: &nav.Location{
						Lat: 30.333,
						Lon: 122.333,
					},
					Distance: 6,
					Duration: 6,
				},
				{
					StationID: "station7",
					StationLocation: &nav.Location{
						Lat: -10.333,
						Lon: 122.333,
					},
					Distance: 7,
					Duration: 7,
				},
			},
			&nav.Location{
				Lat: 32.333,
				Lon: -122.333,
			},
		},
		{
			"station4",
			[]*connectivitymap.QueryResult{
				{
					StationID: "station6",
					StationLocation: &nav.Location{
						Lat: 30.333,
						Lon: 122.333,
					},
					Distance: 8,
					Duration: 8,
				},
				{
					StationID: "station7",
					StationLocation: &nav.Location{
						Lat: -10.333,
						Lon: 122.333,
					},
					Distance: 9,
					Duration: 9,
				},
			},
			&nav.Location{
				Lat: -32.333,
				Lon: 122.333,
			},
		},
		{
			"station6",
			[]*connectivitymap.QueryResult{
				{
					StationID: "dest_location",
					StationLocation: &nav.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
					Distance: 66.6,
					Duration: 66.6,
				},
			},
			&nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
		},
		{
			"station7",
			[]*connectivitymap.QueryResult{
				{
					StationID: "dest_location",
					StationLocation: &nav.Location{
						Lat: 4.4,
						Lon: 4.4,
					},
					Distance: 11.1,
					Duration: 11.1,
				},
			},
			&nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
		},
		{
			"dest_location",
			nil,
			&nav.Location{
				Lat: 4.4,
				Lon: 4.4,
			},
		},
	}

	c := generateNeighborsChan()

	q := NewQuerierBasedOnWeightBetweenNeighborsChan(c)

	for _, c := range cases {
		acturalQueryResult := q.NearByStationQuery(c.stationID)
		if !reflect.DeepEqual(acturalQueryResult, c.expectQueryResult) {
			t.Errorf("Generate incorrect for NearByStationQuery, expect \n%#v\n but got \n%#v\n", c.expectQueryResult, acturalQueryResult)
		}

		actualLocation := q.GetLocation(c.stationID)
		if !reflect.DeepEqual(actualLocation, c.expectLocation) {
			t.Errorf("Generate incorrect for GetLocation, expect \n%#v\n but got \n%#v\n", c.expectLocation, actualLocation)
		}

	}
}

func generateNeighborsChan() chan stationfindertype.WeightBetweenNeighbors {
	c := make(chan stationfindertype.WeightBetweenNeighbors, 3)

	go func() {
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: stationfindertype.NeighborInfoArray0,
			Err:           nil,
		}
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: stationfindertype.NeighborInfoArray1,
			Err:           nil,
		}
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: stationfindertype.NeighborInfoArray2,
			Err:           nil,
		}

		close(c)
	}()

	return c
}
