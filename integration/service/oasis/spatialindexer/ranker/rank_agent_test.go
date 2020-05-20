package ranker

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

func TestRankAgent(t *testing.T) {
	cases := []struct {
		input  []*common.RankedPlaceInfo
		expect []*common.RankedPlaceInfo
	}{
		{
			input: []*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 3,
						Location: &nav.Location{
							Lat: 3.3,
							Lon: 3.3,
						},
					},
					Weight: &common.Weight{
						Distance: 3.3,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 1,
						Location: &nav.Location{
							Lat: 1.1,
							Lon: 1.1,
						},
					},
					Weight: &common.Weight{
						Distance: 1.1,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 22,
						Location: &nav.Location{
							Lat: 22.22,
							Lon: 22.22,
						},
					},
					Weight: &common.Weight{
						Distance: 22.22,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 4,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &common.Weight{
						Distance: 4.4,
					},
				},
			},
			expect: []*common.RankedPlaceInfo{
				{
					PlaceInfo: common.PlaceInfo{
						ID: 1,
						Location: &nav.Location{
							Lat: 1.1,
							Lon: 1.1,
						},
					},
					Weight: &common.Weight{
						Distance: 1.1,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 3,
						Location: &nav.Location{
							Lat: 3.3,
							Lon: 3.3,
						},
					},
					Weight: &common.Weight{
						Distance: 3.3,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 4,
						Location: &nav.Location{
							Lat: 4.4,
							Lon: 4.4,
						},
					},
					Weight: &common.Weight{
						Distance: 4.4,
					},
				},
				{
					PlaceInfo: common.PlaceInfo{
						ID: 22,
						Location: &nav.Location{
							Lat: 22.22,
							Lon: 22.22,
						},
					},
					Weight: &common.Weight{
						Distance: 22.22,
					},
				},
			},
		},
	}

	for _, c := range cases {
		var wg sync.WaitGroup
		pointWithDistanceC := make(chan *common.RankedPlaceInfo, len(c.input))
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for _, item := range c.input {
				pointWithDistanceC <- item
			}
			close(pointWithDistanceC)
		}(&wg)

		wg.Wait()
		rankAgent := newRankAgent(len(c.input))
		actual := rankAgent.RankByDistance(pointWithDistanceC)
		if !reflect.DeepEqual(actual, c.expect) {
			t.Errorf("During test rankAgent, while handling \n input\n %s,\n expect\n %s \n but actual is\n %s\n",
				printRankedPointInfoArray(c.input),
				printRankedPointInfoArray(c.expect),
				printRankedPointInfoArray(actual))
		}
	}
}

func printRankedPointInfoArray(arr []*common.RankedPlaceInfo) string {
	var str string
	for _, item := range arr {
		str += fmt.Sprintf("%#v ", item)
	}
	return str
}
