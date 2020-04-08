package ranker

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

type osrmRanker struct {
	oc *osrmconnector.OSRMConnector
}

func newOsrmRanker(oc *osrmconnector.OSRMConnector) *osrmRanker {
	return &osrmRanker{
		oc: oc,
	}
}

func (ranker *osrmRanker) RankPointIDsByGreatCircleDistance(center spatialindexer.Location, targets []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, targets)
}

func (ranker *osrmRanker) RankPointIDsByShortestDistance(center spatialindexer.Location, targets []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByOSRMShortestPath(center, targets, ranker.oc, pointsThresholdPerRequest)
}
