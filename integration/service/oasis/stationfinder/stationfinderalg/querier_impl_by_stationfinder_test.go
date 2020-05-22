package stationfinderalg

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
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
		placeID           common.PlaceID
		expectQueryResult []*common.RankedPlaceInfo
		expectLocation    *nav.Location
	}{
		{
			stationfindertype.OrigLocationID,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 1,
						Location: &nav.Location{
							Lat: 32.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 22.2,
						Duration: 22.2,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 2,
						Location: &nav.Location{
							Lat: -32.333,
							Lon: -122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 11.1,
						Duration: 11.1,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 3,
						Location: &nav.Location{
							Lat: 32.333,
							Lon: -122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 33.3,
						Duration: 33.3,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 4,
						Location: &nav.Location{
							Lat: -32.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 44.4,
						Duration: 44.4,
					},
				},
			},
			&nav.Location{
				Lat: 1.1,
				Lon: 1.1,
			},
		},
		{
			1,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 2,
						Duration: 2,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 3,
						Duration: 3,
					},
				},
			},
			&nav.Location{
				Lat: 32.333,
				Lon: 122.333,
			},
		},
		{
			2,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 4,
						Duration: 4,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 5,
						Duration: 5,
					},
				},
			},
			&nav.Location{
				Lat: -32.333,
				Lon: -122.333,
			},
		},
		{
			3,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 6,
						Duration: 6,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 7,
						Duration: 7,
					},
				},
			},
			&nav.Location{
				Lat: 32.333,
				Lon: -122.333,
			},
		},
		{
			4,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 8,
						Duration: 8,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &common.Weight{
						Distance: 9,
						Duration: 9,
					},
				},
			},
			&nav.Location{
				Lat: -32.333,
				Lon: 122.333,
			},
		},
		{
			6,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: stationfindertype.DestLocationID,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &common.Weight{
						Distance: 66.6,
						Duration: 66.6,
					},
				},
			},
			&nav.Location{
				Lat: 30.333,
				Lon: 122.333,
			},
		},
		{
			7,
			[]*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: stationfindertype.DestLocationID,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &common.Weight{
						Distance: 11.1,
						Duration: 11.1,
					},
				},
			},
			&nav.Location{
				Lat: -10.333,
				Lon: 122.333,
			},
		},
		{
			stationfindertype.DestLocationID,
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
		acturalQueryResult := q.NearByStationQuery(c.placeID)
		if !reflect.DeepEqual(acturalQueryResult, c.expectQueryResult) {
			t.Errorf("Generate incorrect for NearByStationQuery, expect \n%#v\n but got \n%#v\n", c.expectQueryResult, acturalQueryResult)
		}

		actualLocation := q.GetLocation(c.placeID)
		if !reflect.DeepEqual(actualLocation, c.expectLocation) {
			t.Errorf("Generate incorrect for GetLocation, expect \n%#v\n but got \n%#v\n", c.expectLocation, actualLocation)
		}

	}
}

func generateNeighborsChan() chan stationfindertype.WeightBetweenNeighbors {
	c := make(chan stationfindertype.WeightBetweenNeighbors, 3)

	go func() {
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray0,
			Err:           nil,
		}
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray1,
			Err:           nil,
		}
		c <- stationfindertype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray2,
			Err:           nil,
		}

		close(c)
	}()

	return c
}
