package ranker

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/blevesearch/bleve/geo"
	"github.com/golang/glog"
)

func rankPointsByGreatCircleDistanceToCenter(center spatialindexer.Location, targets []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	if len(targets) == 0 {
		glog.Warningf("When try to rankPointsByGreatCircleDistanceToCenter, input array is empty, center = %+v\n", center)
		return nil
	}

	pointWithDistanceC := make(chan *spatialindexer.RankedPointInfo, len(targets))
	go func() {
		defer close(pointWithDistanceC)

		for _, p := range targets {
			pointWithDistanceC <- &spatialindexer.RankedPointInfo{
				PointInfo: spatialindexer.PointInfo{
					ID:       p.ID,
					Location: p.Location,
				},
				// geo.Haversin's unit is kilometer, convert to meter
				Distance: geo.Haversin(center.Lon, center.Lat, p.Location.Lon, p.Location.Lat) * 1000,
			}
		}
	}()

	rankAgent := newRankAgent(len(targets))
	return rankAgent.RankByDistance(pointWithDistanceC)
}
