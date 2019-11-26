package rankingservice

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/graph"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
	"github.com/golang/glog"
)

func (h *Handler) updateRoutesByTraffic(routes []*osrmv1.Route) []*osrmv1.Route {
	updatedRoutes := []*osrmv1.Route{}
	for _, r := range routes {
		h.updateRouteByTraffic(r)
		if math.IsInf(r.Duration, 0) || math.IsInf(r.Weight, 0) {
			continue // ignore the route
		}
		updatedRoutes = append(updatedRoutes, r)
	}
	return updatedRoutes
}

func (h *Handler) updateRouteByTraffic(route *osrmv1.Route) {
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

func (h *Handler) updateLegByTraffic(leg *osrmv1.RouteLeg) {
	if leg == nil {
		glog.Error("empty leg")
		return
	}
	edges := nodesToEdges(leg.Annotation.Nodes)
	edgesCount := len(edges)
	if len(leg.Annotation.Distance) != edgesCount ||
		len(leg.Annotation.Duration) != edgesCount ||
		len(leg.Annotation.Speed) != edgesCount ||
		len(leg.Annotation.Weight) != edgesCount ||
		len(leg.Annotation.DataSources) != edgesCount {
		glog.Errorf("annotation counts not match")
		return
	}

	if blocked, _ := h.trafficInquirer.EdgesBlockedByIncidents(edges); blocked {
		glog.Warningf("leg blocked by incident, set duration to infinity")
		leg.Duration = math.Inf(0)
		leg.Weight = math.Inf(0)
		return
	}
	flows := h.trafficInquirer.QueryFlows(edges)
	if len(flows) != edgesCount {
		glog.Fatalf("query flow return count %d doesn't match edges count %d", len(flows), edgesCount)
		return
	}

	validFlowsCount := 0
	trafficDataSourceNameIndex := len(leg.Annotation.Metadata.DataSourceNames)
	var sumOriginalAnnotationDuration, sumOriginalAnnotationDistance, sumOriginalAnnotationWeight float64
	var newLegDuration, newLegWeight float64
	for i := range flows {
		sumOriginalAnnotationDistance += leg.Annotation.Distance[i]
		sumOriginalAnnotationDuration += leg.Annotation.Duration[i]
		sumOriginalAnnotationWeight += leg.Annotation.Weight[i]

		if flows[i] != nil {
			if flows[i].IsBlocking() {
				glog.Warningf("leg blocked by flow, set duration to infinity")
				leg.Duration = math.Inf(0)
				leg.Weight = math.Inf(0)
				return // exit the update once found blocking flow
			}

			validFlowsCount++
			leg.Annotation.DataSources[i] = trafficDataSourceNameIndex
			leg.Annotation.Speed[i] = float64(flows[i].Speed)
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
	glog.V(2).Infof("leg edges count %d, leg distance %f(%f sum by annotation), valid flows count %d, duration,weight %f,%f(%f,%f sum by annotation, %f,%f on turn) -> %f,%f (reduced %f,%f)",
		edgesCount, leg.Distance, sumOriginalAnnotationDistance, validFlowsCount,
		leg.Duration, leg.Weight, sumOriginalAnnotationDuration, sumOriginalAnnotationWeight, legTurnDuration, legTurnWeight, newLegDuration, newLegWeight, leg.Duration-newLegDuration, leg.Weight-newLegWeight)

	if validFlowsCount > 0 {
		leg.Annotation.Metadata.DataSourceNames = append(leg.Annotation.Metadata.DataSourceNames, "traffic cache")
		leg.Duration = newLegDuration
		leg.Weight = newLegWeight
	}
}

func nodesToEdges(nodes []int64) []graph.Edge {
	edges := []graph.Edge{}

	for i := 0; i < len(nodes)-1; i++ {
		edges = append(edges, graph.Edge{From: nodes[i], To: nodes[i+1]})
	}

	return edges
}
