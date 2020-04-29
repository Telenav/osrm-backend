package ranking

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/golang/glog"
)

func (h *Handler) updateRoutesByTraffic(routes []*route.Route) []*route.Route {
	updatedRoutes := []*route.Route{}
	for _, r := range routes {
		h.updateRouteByTraffic(r)
		if math.IsInf(r.Duration, 0) || math.IsInf(r.Weight, 0) {
			continue // ignore the route
		}
		updatedRoutes = append(updatedRoutes, r)
	}
	return updatedRoutes
}

func (h *Handler) updateRouteByTraffic(route *route.Route) {
	if route == nil {
		glog.Error("empty route")
		return
	}

	var newRouteDuration, newRouteWeight float64
	for _, l := range route.Legs {
		h.updateLegByTraffic(l)
		if math.IsInf(l.Duration, 0) || math.IsInf(l.Weight, 0) {
			glog.Warningf("route blocked by live traffic, set duration to infinity")
			route.Duration = l.Duration
			route.Weight = l.Weight
			return
		}
		newRouteDuration += l.Duration
		newRouteWeight += l.Weight
	}
	glog.V(1).Infof("route update by traffic, duration,weight %f,%f -> %f,%f (reduced %f,%f)",
		route.Duration, route.Weight, newRouteDuration, newRouteWeight, route.Duration-newRouteDuration, route.Weight-newRouteWeight)
	route.Duration = newRouteDuration
	route.Weight = newRouteWeight
}

func (h *Handler) updateLegByTraffic(leg *route.Leg) {
	if leg == nil {
		glog.Error("empty leg")
		return
	}
	waysCount := len(leg.Annotation.Ways)
	if len(leg.Annotation.Distance) != waysCount ||
		len(leg.Annotation.Duration) != waysCount ||
		len(leg.Annotation.Speed) != waysCount ||
		len(leg.Annotation.Weight) != waysCount ||
		len(leg.Annotation.DataSources) != waysCount {
		glog.Errorf("annotation counts not match")
		return
	}

	validFlowsCount := 0
	trafficDataSourceNameIndex := len(leg.Annotation.Metadata.DataSourceNames)
	var sumOriginalAnnotationDuration, sumOriginalAnnotationDistance, sumOriginalAnnotationWeight float64
	var newLegDuration, newLegWeight float64
	for i := 0; i < waysCount; i++ {
		sumOriginalAnnotationDistance += leg.Annotation.Distance[i]
		sumOriginalAnnotationDuration += leg.Annotation.Duration[i]
		sumOriginalAnnotationWeight += leg.Annotation.Weight[i]

		wayID := leg.Annotation.Ways[i]

		if h.trafficQuerier.BlockedByIncident(wayID) {
			glog.Warningf("way %d on leg blocked by incident, set duration to infinity", wayID)
			leg.Duration = math.Inf(0)
			leg.Weight = math.Inf(0)
			return
		}

		flow := h.trafficQuerier.QueryFlow(wayID)
		if flow != nil {
			if flow.IsBlocking() {
				glog.Warningf("way %d on leg blocked by flow %v, set duration to infinity", wayID, flow)
				leg.Duration = math.Inf(0)
				leg.Weight = math.Inf(0)
				return // exit the update once found blocking flow
			}

			validFlowsCount++
			leg.Annotation.DataSources[i] = trafficDataSourceNameIndex
			leg.Annotation.Speed[i] = float64(flow.Speed)
			leg.Annotation.Duration[i] = leg.Annotation.Distance[i] / leg.Annotation.Speed[i] // not include turn duration
			leg.Annotation.Weight[i] = leg.Annotation.Distance[i] / leg.Annotation.Speed[i]   // not include turn weight
		}
		newLegDuration += leg.Annotation.Duration[i]
		newLegWeight += leg.Annotation.Weight[i]
	}
	// metioned in doc, duration and weight in annotation does not include any on turn, so we calculate turn duration/weight by minus.
	// https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#annotation-object
	legTurnDuration := leg.Duration - sumOriginalAnnotationDuration
	legTurnWeight := leg.Weight - sumOriginalAnnotationWeight
	newLegDuration += legTurnDuration
	newLegWeight += legTurnWeight
	glog.V(2).Infof("leg ways count %d, leg distance %f(%f sum by annotation), valid flows count %d, duration,weight %f,%f(%f,%f sum by annotation, %f,%f on turn) -> %f,%f (reduced %f,%f)",
		waysCount, leg.Distance, sumOriginalAnnotationDistance, validFlowsCount,
		leg.Duration, leg.Weight, sumOriginalAnnotationDuration, sumOriginalAnnotationWeight, legTurnDuration, legTurnWeight, newLegDuration, newLegWeight, leg.Duration-newLegDuration, leg.Weight-newLegWeight)

	if validFlowsCount > 0 {
		leg.Annotation.Metadata.DataSourceNames = append(leg.Annotation.Metadata.DataSourceNames, "live traffic")
		leg.Duration = newLegDuration
		leg.Weight = newLegWeight
	}
}
