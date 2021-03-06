package ranker

import (
	"sort"

	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

// rankAgent accepts items to be ranked then returns ranking result
type rankAgent struct {
	rankedPoints []*entity.TransferInfo
}

func newRankAgent(pointNum int) *rankAgent {
	return &rankAgent{
		rankedPoints: make([]*entity.TransferInfo, 0, pointNum),
	}
}

type rankItems []*entity.TransferInfo

func (r rankItems) Len() int {
	return len(r)
}

func (r rankItems) Less(i, j int) bool {
	return r[i].Weight.Distance < r[j].Weight.Distance
}

func (r rankItems) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *rankAgent) RankByDistance(input <-chan *entity.TransferInfo) []*entity.TransferInfo {
	for p := range input {
		r.rankedPoints = append(r.rankedPoints, p)
	}

	sort.Sort(rankItems(r.rankedPoints))

	return r.rankedPoints
}
