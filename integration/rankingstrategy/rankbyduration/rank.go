package rankbyduration

import (
	"sort"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
)

type rankItems []*osrmv1.Route

func (r rankItems) Len() int {
	return len(r)
}

func (r rankItems) Less(i, j int) bool {
	return r[i].Duration < r[j].Duration
}

func (r rankItems) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Rank ranks routes by duration.
func Rank(routes []*osrmv1.Route) []*osrmv1.Route {
	sort.Sort(rankItems(routes))
	return routes
}
