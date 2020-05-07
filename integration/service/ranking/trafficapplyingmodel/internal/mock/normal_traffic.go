package mock

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
)

// NormalTraffic implements both traffic.HistoricalSpeedQuerier and traffic.LiveTrafficQuerier for the OSRMRouteNormal.
type NormalTraffic struct {
	blockedWays      map[int64]struct{}
	flows            map[int64]*trafficproxy.Flow
	historicalspeeds map[int64]float64
}

// NewNormalTraffic creates new NormalTraffic object.
func NewNormalTraffic() NormalTraffic {
	n := NewNormalTrafficNoBlock()
	n.blockedWays[888079015] = struct{}{}
	return n
}

// NewNormalTrafficNoBlock creates new NormalTraffic object without any block way.
func NewNormalTrafficNoBlock() NormalTraffic {
	n := NormalTraffic{
		map[int64]struct{}{}, map[int64]*trafficproxy.Flow{}, map[int64]float64{},
	}

	n.flows[851242314] = &trafficproxy.Flow{WayID: 851242314, Speed: 6.110000, TrafficLevel: trafficproxy.TrafficLevel_SLOW_SPEED, Timestamp: 1579419488000}
	n.flows[-23704643] = &trafficproxy.Flow{WayID: -23704643, Speed: 106.11, TrafficLevel: trafficproxy.TrafficLevel_FREE_FLOW, Timestamp: 1579419488000}
	n.historicalspeeds[1234704366] = 20.5
	n.historicalspeeds[-23704642] = 70.0
	return n
}

// BlockedByIncident implements BlockedByIncident of traffic.LiveTrafficQuerier but always return false.
func (n NormalTraffic) BlockedByIncident(wayID int64) bool {
	_, ok := n.blockedWays[wayID]
	return ok
}

// QueryFlow implements QueryFlow of traffic.LiveTrafficQuerier but always return nil.
func (n NormalTraffic) QueryFlow(wayID int64) *trafficproxy.Flow {
	if f, ok := n.flows[wayID]; ok {
		return f
	}
	return nil
}

// QueryHistoricalSpeed implements QueryHistoricalSpeed of traffic.HistoricalSpeedQuerier but always return false.
func (n NormalTraffic) QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool) {
	if s, ok := n.historicalspeeds[wayID]; ok {
		return s, true
	}
	return 0, false
}
