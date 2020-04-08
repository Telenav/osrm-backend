package ranker

import "github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"

type simpleRanker struct {
}

func newSimpleRanker() *simpleRanker {
	return &simpleRanker{}
}

func (ranker *simpleRanker) RankPointIDsByGreatCircleDistance(center spatialindexer.Location,
	targets []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return rankPointsByGreatCircleDistanceToCenter(center, targets)
}

func (ranker *simpleRanker) RankPointIDsByShortestDistance(center spatialindexer.Location,
	targets []*spatialindexer.PointInfo) []*spatialindexer.RankedPointInfo {
	return ranker.RankPointIDsByGreatCircleDistance(center, targets)
}
