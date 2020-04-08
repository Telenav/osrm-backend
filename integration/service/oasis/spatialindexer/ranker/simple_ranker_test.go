package ranker

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

func TestRankerInterfaceViaSimpleRanker(t *testing.T) {
	cases := []struct {
		center  spatialindexer.Location
		targets []*spatialindexer.PointInfo
		expect  []*spatialindexer.RankedPointInfo
	}{
		{
			center: spatialindexer.Location{
				Lat: 37.398973,
				Lon: -121.976633,
			},
			targets: []*spatialindexer.PointInfo{
				&spatialindexer.PointInfo{
					ID: 1,
					Location: spatialindexer.Location{
						Lat: 37.388840,
						Lon: -121.981736,
					},
				},
				&spatialindexer.PointInfo{
					ID: 2,
					Location: spatialindexer.Location{
						Lat: 37.375515,
						Lon: -121.942812,
					},
				},
				&spatialindexer.PointInfo{
					ID: 3,
					Location: spatialindexer.Location{
						Lat: 37.336954,
						Lon: -121.861624,
					},
				},
			},
			expect: []*spatialindexer.RankedPointInfo{
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 1,
						Location: spatialindexer.Location{
							Lat: 37.388840,
							Lon: -121.981736,
						},
					},
					Distance: 1213.445757354474,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 2,
						Location: spatialindexer.Location{
							Lat: 37.375515,
							Lon: -121.942812,
						},
					},
					Distance: 3965.986474110687,
				},
				&spatialindexer.RankedPointInfo{
					PointInfo: spatialindexer.PointInfo{
						ID: 3,
						Location: spatialindexer.Location{
							Lat: 37.336954,
							Lon: -121.861624,
						},
					},
					Distance: 12281.070927352637,
				},
			},
		},
	}

	ranker := CreateRanker(SimpleRanker, nil)

	for _, c := range cases {
		actual := ranker.RankPointIDsByGreatCircleDistance(c.center, c.targets)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During test SimpleRanker's RankPointIDsByGreatCircleDistance, \n expect \n%s \nwhile actual is\n %s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}

		actual = ranker.RankPointIDsByShortestDistance(c.center, c.targets)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During test SimpleRanker's RankPointIDsByGreatCircleDistance, \n expect \n%s \nwhile actual is\n %s\n",
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}
	}

}
