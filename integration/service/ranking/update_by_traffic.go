package ranking

import (
	"math"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
)

func (h *Handler) updateRoutesByTraffic(routes []*route.Route, enableLiveTraffic, enableHistoricalSpeed bool) ([]*route.Route, error) {
	updatedRoutes := []*route.Route{}
	for _, r := range routes {

		if err := h.trafficApplier.ApplyTraffic(r, enableLiveTraffic, enableHistoricalSpeed); err != nil {
			return nil, err
		}
		if math.IsInf(r.Duration, 0) || math.IsInf(r.Weight, 0) {
			continue // discard the route
		}
		updatedRoutes = append(updatedRoutes, r)
	}
	return updatedRoutes, nil
}

func parseTrafficOptions(liveTrafficQueryValue, historicalSpeedQueryValue string) (bool, bool) {
	var enableLiveTraffic, enableHistoricalSpeed bool // both false by default
	if b, err := strconv.ParseBool(liveTrafficQueryValue); err == nil {
		enableLiveTraffic = b
	}
	if b, err := strconv.ParseBool(historicalSpeedQueryValue); err == nil {
		enableHistoricalSpeed = b
	}
	return enableLiveTraffic, enableHistoricalSpeed
}

const (
	// NOTE: These two are route request keys, but they're not OSRM original keys.
	// We keep them here to reduce impacts but provides flexibility to support these options on the fly.
	// Consider to move them into OSRM API when really necessary.
	liveTrafficQueryKey     = "live_traffic"
	historicalSpeedQueryKey = "historical_speed"
)
