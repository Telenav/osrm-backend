// Package appendspeedonly only additional live traffic and historical speed append to OSRM response, keep original response's duration/weight for a route.
// It queries live traffic and historical speed at current clock and append to annotations on a OSRM route.
// The live/historical speed will be set < 0 if no valid data.
package appendspeedonly

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/traffic"
)

// Name represents the name of traffic applying model.
func Name() string {
	return "append-speed-only"
}

// Model represents the model of applying traffic.
type Model struct {
	traffic.LiveTrafficQuerier
	traffic.HistoricalSpeedQuerier
}

// New creates a new model object.
func New(l traffic.LiveTrafficQuerier, h traffic.HistoricalSpeedQuerier) (*Model, error) {
	if l == nil && h == nil {
		return nil, trafficapplyingmodel.ErrNoValidTrafficQuerier
	}

	return &Model{l, h}, nil
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

	for _, l := range r.Legs {
		if l == nil {
			return trafficapplyingmodel.ErrEmptyLeg
		}
		if err := m.applyTrafficOnAnnotation(l.Annotation, liveTraffic, historicalSpeed); err != nil {
			return err
		}
	}

	return nil
}

func (m Model) applyTrafficOnAnnotation(a *route.Annotation, enableLiveTraffic bool, enableHistoricalSpeed bool) error {
	if a == nil {
		return trafficapplyingmodel.ErrEmptyAnnotation
	}

	waysCount := len(a.Ways)

	liveTrafficSpeeds := make([]float64, waysCount)
	liveTrafficlevels := make([]int, waysCount)
	blockIncidents := make([]bool, waysCount)
	historicalSpeeds := make([]float64, waysCount)

	utcTimestamp := time.Now().UTC()
	for i := 0; i < waysCount; i++ {
		wayID := a.Ways[i]

		if m.LiveTrafficQuerier != nil && enableLiveTraffic {
			if f := m.LiveTrafficQuerier.QueryFlow(wayID); f != nil {
				liveTrafficSpeeds[i] = float64(f.GetSpeed())
				liveTrafficlevels[i] = int(f.GetTrafficLevel())
			} else {
				liveTrafficSpeeds[i] = trafficapplyingmodel.InvalidTrafficSpeedFloat64
				liveTrafficlevels[i] = int(trafficproxy.TrafficLevel_NO_LEVELS)
			}
			if m.LiveTrafficQuerier.BlockedByIncident(wayID) {
				blockIncidents[i] = true
			}
		}

		if m.HistoricalSpeedQuerier != nil && enableHistoricalSpeed {
			if speed, ok := m.HistoricalSpeedQuerier.QueryHistoricalSpeed(wayID, utcTimestamp); ok {
				historicalSpeeds[i] = speed
			} else {
				historicalSpeeds[i] = trafficapplyingmodel.InvalidTrafficSpeedFloat64
			}
		}
	}

	if m.LiveTrafficQuerier != nil && enableLiveTraffic {
		a.LiveTrafficLevel = liveTrafficlevels
		a.LiveTrafficSpeed = liveTrafficSpeeds
		a.BlockIncident = blockIncidents
	}
	if m.HistoricalSpeedQuerier != nil && enableHistoricalSpeed {
		a.HistoricalSpeed = historicalSpeeds
	}
	return nil
}
