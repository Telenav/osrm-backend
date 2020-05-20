package ranker

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"
)

type simpleRanker struct {
}

func newSimpleRanker() *simpleRanker {
	return &simpleRanker{}
}

func (ranker *simpleRanker) RankPlaceIDsByGreatCircleDistance(center nav.Location,
	targets []*common.PlaceInfo) []*common.RankedPlaceInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, targets)
}

func (ranker *simpleRanker) RankPlaceIDsByShortestDistance(center nav.Location,
	targets []*common.PlaceInfo) []*common.RankedPlaceInfo {
	return ranker.RankPlaceIDsByGreatCircleDistance(center, targets)
}
