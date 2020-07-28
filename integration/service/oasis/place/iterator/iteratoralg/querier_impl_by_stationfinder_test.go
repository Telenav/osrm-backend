package iteratoralg

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
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
		placeID           entity.PlaceID
		expectQueryResult []*entity.TransferInfo
		expectLocation    *nav.Location
	}{
		{
			iteratortype.OrigLocationID,
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 1,
						Location: &nav.Location{
							Lat: 32.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 22.2,
						Duration: 22.2,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 2,
						Location: &nav.Location{
							Lat: -32.333,
							Lon: -122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 11.1,
						Duration: 11.1,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 3,
						Location: &nav.Location{
							Lat: 32.333,
							Lon: -122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 33.3,
						Duration: 33.3,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 4,
						Location: &nav.Location{
							Lat: -32.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 2,
						Duration: 2,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 4,
						Duration: 4,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 6,
						Duration: 6,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 6,
						Location: &nav.Location{
							Lat: 30.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
						Distance: 8,
						Duration: 8,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 7,
						Location: &nav.Location{
							Lat: -10.333,
							Lon: 122.333,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: iteratortype.DestLocationID,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &entity.Weight{
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
			[]*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: iteratortype.DestLocationID,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &entity.Weight{
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
			iteratortype.DestLocationID,
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
		acturalQueryResult := q.GetConnectedPlaces(c.placeID)
		if !reflect.DeepEqual(acturalQueryResult, c.expectQueryResult) {
			t.Errorf("Generate incorrect for GetConnectedPlaces, expect \n%#v\n but got \n%#v\n", c.expectQueryResult, acturalQueryResult)
		}

		actualLocation := q.GetLocation(c.placeID)
		if !reflect.DeepEqual(actualLocation, c.expectLocation) {
			t.Errorf("Generate incorrect for GetLocation, expect \n%#v\n but got \n%#v\n", c.expectLocation, actualLocation)
		}

	}
}

func generateNeighborsChan() chan iteratortype.WeightBetweenNeighbors {
	c := make(chan iteratortype.WeightBetweenNeighbors, 3)

	go func() {
		c <- iteratortype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray0,
			Err:           nil,
		}
		c <- iteratortype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray1,
			Err:           nil,
		}
		c <- iteratortype.WeightBetweenNeighbors{
			NeighborsInfo: mock.NeighborInfoArray2,
			Err:           nil,
		}

		close(c)
	}()

	return c
}
