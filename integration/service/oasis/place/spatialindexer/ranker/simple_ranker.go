package ranker

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

type simpleRanker struct {
}

func newSimpleRanker() *simpleRanker {
	return &simpleRanker{}
}

func (ranker *simpleRanker) RankPlaceIDsByGreatCircleDistance(center nav.Location,
	targets []*entity.PlaceWithLocation) []*entity.TransferInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, targets)
}

func (ranker *simpleRanker) RankPlaceIDsByShortestDistance(center nav.Location,
	targets []*entity.PlaceWithLocation) []*entity.TransferInfo {
	return ranker.RankPlaceIDsByGreatCircleDistance(center, targets)
}
