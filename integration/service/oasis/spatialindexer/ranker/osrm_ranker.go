package ranker

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
)

type osrmRanker struct {
	oc *osrmconnector.OSRMConnector
}

func newOsrmRanker(oc *osrmconnector.OSRMConnector) *osrmRanker {
	return &osrmRanker{
		oc: oc,
	}
}

func (ranker *osrmRanker) RankPlaceIDsByGreatCircleDistance(center nav.Location, targets []*common.PlaceInfo) []*common.RankedPlaceInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, targets)
}

func (ranker *osrmRanker) RankPlaceIDsByShortestDistance(center nav.Location, targets []*common.PlaceInfo) []*common.RankedPlaceInfo {
	return rankPointsByOSRMShortestPath(center, targets, ranker.oc, pointsThresholdPerRequest)
}
