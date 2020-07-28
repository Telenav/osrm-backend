package ranker

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/blevesearch/bleve/geo"
	"github.com/golang/glog"
)

// defaultSpeed is 22.2 meter/second
// 22.2 m/s = 50 miles/hour
const defaultSpeed = 22.2

func rankPointsByGreatCircleDistanceToCenter(center nav.Location, targets []*entity.PlaceWithLocation) []*entity.TransferInfo {
	if len(targets) == 0 {
		glog.Warningf("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty, center = %+v\n", center)
		return nil
	}

	pointWithDistanceC := make(chan *entity.TransferInfo, len(targets))
	go func() {
		defer close(pointWithDistanceC)

		for _, p := range targets {
			// geo.Haversin's unit is kilometer, convert to meter
			length := geo.Haversin(center.Lon, center.Lat, p.Location.Lon, p.Location.Lat) * 1000

			pointWithDistanceC <- &entity.TransferInfo{
				PlaceWithLocation: entity.PlaceWithLocation{
					ID:       p.ID,
					Location: p.Location,
				},
				Weight: &entity.Weight{
					Distance: length,
					Duration: length / defaultSpeed,
				},
			}
		}
	}()

	rankAgent := newRankAgent(len(targets))
	return rankAgent.RankByDistance(pointWithDistanceC)
}
