package ranker

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

func TestRankerInterfaceViaSimpleRanker(t *testing.T) {
	cases := []struct {
		center  nav.Location
		targets []*entity.PlaceWithLocation
		expect  []*entity.TransferInfo
	}{
		{
			center: nav.Location{
				Lat: 37.398973,
				Lon: -121.976633,
			},
			targets: []*entity.PlaceWithLocation{
				{
					ID: 1,
					Location: &nav.Location{
						Lat: 37.388840,
						Lon: -121.981736,
					},
				},
				{
					ID: 2,
					Location: &nav.Location{
						Lat: 37.375515,
						Lon: -121.942812,
					},
				},
				{
					ID: 3,
					Location: &nav.Location{
						Lat: 37.336954,
						Lon: -121.861624,
					},
				},
			},
			expect: []*entity.TransferInfo{
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 1,
						Location: &nav.Location{
							Lat: 37.388840,
							Lon: -121.981736,
						},
					},
					Weight: &entity.Weight{
						Distance: 1213.445757354474,
						Duration: 54.65971879975108,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 2,
						Location: &nav.Location{
							Lat: 37.375515,
							Lon: -121.942812,
						},
					},
					Weight: &entity.Weight{
						Distance: 3965.986474110687,
						Duration: 178.64803937435528,
					},
				},
				{
					PlaceWithLocation: entity.PlaceWithLocation{
						ID: 3,
						Location: &nav.Location{
							Lat: 37.336954,
							Lon: -121.861624,
						},
					},
					Weight: &entity.Weight{
						Distance: 12281.070927352637,
						Duration: 553.2013931239927,
					},
				},
			},
		},
	}

	ranker := CreateRanker(SimpleRanker, nil)

	for _, c := range cases {
		actual := ranker.RankPlaceIDsByGreatCircleDistance(c.center, c.targets)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During test SimpleRanker's RankPlaceIDsByGreatCircleDistance, \n expect \n%s \nwhile actual is\n %s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}

		actual = ranker.RankPlaceIDsByShortestDistance(c.center, c.targets)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During test SimpleRanker's RankPlaceIDsByGreatCircleDistance, \n expect \n%s \nwhile actual is\n %s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}
	}

}
