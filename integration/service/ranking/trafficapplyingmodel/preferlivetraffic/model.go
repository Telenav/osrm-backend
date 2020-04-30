// Package preferlivetraffic re-calculate duration/weight for a route, prefer to use live traffic, if no live traffic on a way then use historical speed, otherwise use the speed from Lua profile.
// Both live traffic and historical speed will be appended to annotations on a OSRM route, same as what package appendspeedonly provides.
package preferlivetraffic

import (
	"fmt"
	"math"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/appendspeedonly"
	"github.com/Telenav/osrm-backend/integration/traffic"
	"github.com/golang/glog"
)

// Name represents the name of traffic applying model.
func Name() string {
	return "prefer-live-traffic"
}

// Model represents the model of applying traffic.
type Model struct {
	traffic.LiveTrafficQuerier
	traffic.HistoricalSpeedQuerier

	appendOnly *appendspeedonly.Model
}

// New creates a new model object.
func New(l traffic.LiveTrafficQuerier, h traffic.HistoricalSpeedQuerier) (*Model, error) {
	if l == nil && h == nil {
		return nil, trafficapplyingmodel.ErrNoValidTrafficQuerier
	}

	appendOnly, err := appendspeedonly.New(l, h)
	if err != nil {
		return nil, err
	}

	return &Model{l, h, appendOnly}, nil
}

// ApplyTraffic applys traffic on a route.
func (m Model) ApplyTraffic(r *route.Route, liveTraffic bool, historicalSpeed bool) error {
	if r == nil {
		return trafficapplyingmodel.ErrEmptyRoute
	}
	if m.LiveTrafficQuerier == nil && m.HistoricalSpeedQuerier == nil {
		return trafficapplyingmodel.ErrNoValidTrafficQuerier
	}
	if !liveTraffic && !historicalSpeed {
		return nil // nothing need to do
	}

	// firstly append traffic at current clock
	if err := m.appendOnly.ApplyTraffic(r, liveTraffic, historicalSpeed); err != nil {
		return err
	}

	// calculate by traffic: prefer live traffic if exist, other historical speed if exist
	var newRouteDuration, newRouteWeight float64
	for _, l := range r.Legs {
		if err := m.recalcOnLeg(l, liveTraffic, historicalSpeed); err != nil {
			return err
		}
		if math.IsInf(l.Duration, 0) || math.IsInf(l.Weight, 0) {
			glog.Warningf("route blocked by live traffic, set duration to infinity")
			r.Duration = l.Duration
			r.Weight = l.Weight
			return nil
		}
		newRouteDuration += l.Duration
		newRouteWeight += l.Weight
	}
	glog.V(1).Infof("route applied traffic, duration,weight %f,%f -> %f,%f (reduced %f,%f)",
		r.Duration, r.Weight, newRouteDuration, newRouteWeight, r.Duration-newRouteDuration, r.Weight-newRouteWeight)
	r.Duration = newRouteDuration
	r.Weight = newRouteWeight
	return nil
}

func (m Model) recalcOnLeg(l *route.Leg, enableLiveTraffic bool, enableHistoricalSpeed bool) error {
	if l == nil {
		return trafficapplyingmodel.ErrEmptyLeg
	}
	if l.Annotation == nil {
		return trafficapplyingmodel.ErrEmptyAnnotation
	}
	if l.Annotation.Metadata == nil {
		return trafficapplyingmodel.ErrEmptyAnnotationMetadata
	}

	// data validating
	waysCount := len(l.Annotation.Ways)
	if len(l.Annotation.Distance) != waysCount ||
		len(l.Annotation.Duration) != waysCount ||
		len(l.Annotation.Speed) != waysCount ||
		len(l.Annotation.Weight) != waysCount ||
		len(l.Annotation.DataSources) != waysCount {
		return fmt.Errorf("annotation data not match: ways,distance,duration,speed,weight,datasources %d,%d,%d,%d,%d,%d",
			waysCount, len(l.Annotation.Distance), len(l.Annotation.Duration), len(l.Annotation.Speed), len(l.Annotation.Weight), len(l.Annotation.DataSources))
	}
	if m.LiveTrafficQuerier != nil && enableLiveTraffic { // requires live traffic data
		if len(l.Annotation.LiveTrafficLevel) != waysCount ||
			len(l.Annotation.LiveTrafficSpeed) != waysCount ||
			len(l.Annotation.BlockIncident) != waysCount {
			return fmt.Errorf("annotation live traffic data not match: ways,live_traffic_level,live_traffic_speed,block_incident %d,%d,%d,%d",
				waysCount, len(l.Annotation.LiveTrafficLevel), len(l.Annotation.LiveTrafficSpeed), len(l.Annotation.BlockIncident))
		}
	}
	if m.HistoricalSpeedQuerier != nil && enableHistoricalSpeed { // requires historical speed data
		if len(l.Annotation.HistoricalSpeed) != waysCount {
			return fmt.Errorf("annotation historical speed data not match: ways,historical_speed %d,%d", waysCount, len(l.Annotation.HistoricalSpeed))
		}
	}

	// set data source index
	var liveTrafficSourceNameIndex, historicalSpeedSourceNameIndex int
	if m.HistoricalSpeedQuerier != nil && enableHistoricalSpeed { // requires historical speed data
		historicalSpeedSourceNameIndex = len(l.Annotation.Metadata.DataSourceNames)
		l.Annotation.Metadata.DataSourceNames = append(l.Annotation.Metadata.DataSourceNames, trafficapplyingmodel.SourceNameHistoricalSpeed)
	}
	if m.LiveTrafficQuerier != nil && enableLiveTraffic { // requires live traffic data
		liveTrafficSourceNameIndex = len(l.Annotation.Metadata.DataSourceNames)
		l.Annotation.Metadata.DataSourceNames = append(l.Annotation.Metadata.DataSourceNames, trafficapplyingmodel.SourceNameLiveTraffic)
	}

	// recalc
	var sumOriginalAnnotationDuration, sumOriginalAnnotationDistance, sumOriginalAnnotationWeight float64
	var newLegDuration, newLegWeight float64
	var appliedHistoricalSpeedCount, appliedLiveTrafficSpeedCount int
	for i := 0; i < waysCount; i++ {
		wayID := l.Annotation.Ways[i]

		sumOriginalAnnotationDistance += l.Annotation.Distance[i]
		sumOriginalAnnotationDuration += l.Annotation.Duration[i]
		sumOriginalAnnotationWeight += l.Annotation.Weight[i]

		if m.HistoricalSpeedQuerier != nil && enableHistoricalSpeed {
			if !trafficapplyingmodel.IsInvalidSpeed(l.Annotation.HistoricalSpeed[i]) { // valid historical speed
				l.Annotation.DataSources[i] = historicalSpeedSourceNameIndex
				l.Annotation.Speed[i] = l.Annotation.HistoricalSpeed[i]
				l.Annotation.Duration[i] = l.Annotation.Distance[i] / l.Annotation.Speed[i] // not include turn duration
				l.Annotation.Weight[i] = l.Annotation.Distance[i] / l.Annotation.Speed[i]   // not include turn weight
			}
		}
		if m.LiveTrafficQuerier != nil && enableLiveTraffic {
			if l.Annotation.BlockIncident[i] {
				glog.Warningf("way %d on leg blocked by incident, set duration to infinity", wayID)
				l.Duration = math.Inf(0)
				l.Weight = math.Inf(0)
				return nil
			}
			if trafficproxy.TrafficLevel(l.Annotation.LiveTrafficLevel[i]) == trafficproxy.TrafficLevel_CLOSED {
				glog.Warningf("way %d on leg blocked by CLOSED flow, set duration to infinity", wayID)
				l.Duration = math.Inf(0)
				l.Weight = math.Inf(0)
				return nil
			}

			if !trafficapplyingmodel.IsInvalidSpeed(l.Annotation.LiveTrafficSpeed[i]) { // valid live traffic speed
				l.Annotation.DataSources[i] = liveTrafficSourceNameIndex
				l.Annotation.Speed[i] = l.Annotation.LiveTrafficSpeed[i]
				l.Annotation.Duration[i] = l.Annotation.Distance[i] / l.Annotation.Speed[i] // not include turn duration
				l.Annotation.Weight[i] = l.Annotation.Distance[i] / l.Annotation.Speed[i]   // not include turn weight
			}
		}

		newLegDuration += l.Annotation.Duration[i]
		newLegWeight += l.Annotation.Weight[i]
		if l.Annotation.DataSources[i] == liveTrafficSourceNameIndex {
			appliedLiveTrafficSpeedCount++
		}
		if l.Annotation.DataSources[i] == historicalSpeedSourceNameIndex {
			appliedHistoricalSpeedCount++
		}
	}
	// metioned in doc, duration and weight in annotation does not include any on turn, so we calculate turn duration/weight by minus.
	// https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#annotation-object
	legTurnDuration := l.Duration - sumOriginalAnnotationDuration
	legTurnWeight := l.Weight - sumOriginalAnnotationWeight
	newLegDuration += legTurnDuration
	newLegWeight += legTurnWeight
	glog.V(2).Infof("leg ways count %d, applied %s %d, appiled %s %d, leg distance %f(%f sum by annotation), duration,weight %f,%f(%f,%f sum by annotation, %f,%f on turn) -> %f,%f (reduced %f,%f)",
		waysCount, trafficapplyingmodel.SourceNameLiveTraffic, appliedLiveTrafficSpeedCount, trafficapplyingmodel.SourceNameHistoricalSpeed, appliedHistoricalSpeedCount,
		l.Distance, sumOriginalAnnotationDistance,
		l.Duration, l.Weight, sumOriginalAnnotationDuration, sumOriginalAnnotationWeight, legTurnDuration, legTurnWeight, newLegDuration, newLegWeight, l.Duration-newLegDuration, l.Weight-newLegWeight)

	if appliedLiveTrafficSpeedCount > 0 || appliedHistoricalSpeedCount > 0 {
		l.Duration = newLegDuration
		l.Weight = newLegWeight
	}

	return nil
}
