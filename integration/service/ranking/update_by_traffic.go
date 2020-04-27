package ranking

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

func (h *Handler) updateRoutesByTraffic(routes []*route.Route) ([]*route.Route, error) {
	updatedRoutes := []*route.Route{}
	for _, r := range routes {

		if err := h.trafficApplier.ApplyTraffic(r, true, true); err != nil {
			return nil, err
		}
		if math.IsInf(r.Duration, 0) || math.IsInf(r.Weight, 0) {
			continue // discard the route
		}
		updatedRoutes = append(updatedRoutes, r)
	}
	return updatedRoutes, nil
}
