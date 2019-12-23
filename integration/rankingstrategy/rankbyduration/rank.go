package rankbyduration

import (
	"sort"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm"
)

type rankItems []*osrm.Route

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
func Rank(routes []*osrm.Route) []*osrm.Route {
	sort.Sort(rankItems(routes))
	return routes
}
