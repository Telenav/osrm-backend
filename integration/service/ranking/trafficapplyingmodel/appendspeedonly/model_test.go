package appendspeedonly_test

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"

	"github.com/Telenav/osrm-backend/integration/traffic"

	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/appendspeedonly"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/internal/mock"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
)

func TestApplyTrafficErrors(t *testing.T) {

	cases := []struct {
		r                      *route.Route
		liveTrafficQuerier     traffic.LiveTrafficQuerier
		historicalSpeedQuerier traffic.HistoricalSpeedQuerier
		liveTraffic            bool
		historicalSpeed        bool
		expect                 error
	}{
		{nil, mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, trafficapplyingmodel.ErrEmptyRoute},
		{mock.NewOSRMRouteNoLeg(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, nil}, // do nothing
		{mock.NewOSRMRouteOneEmptyLeg(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, trafficapplyingmodel.ErrEmptyLeg},
		{mock.NewOSRMRouteOneEmptyLeg(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, false, false, nil}, // do nothing
		{mock.NewOSRMRouteNoAnnotation(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, trafficapplyingmodel.ErrEmptyAnnotation},
		{mock.NewOSRMRouteNormal(), nil, nil, true, true, trafficapplyingmodel.ErrNoValidTrafficQuerier},
		{mock.NewOSRMRouteNormal(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, nil}, // do nothing
	}

	m, err := appendspeedonly.New(mock.EmptyTraffic{}, mock.EmptyTraffic{})
	if err != nil {
		t.Error(err)
	}

	for _, c := range cases {
		m.LiveTrafficQuerier = c.liveTrafficQuerier
		m.HistoricalSpeedQuerier = c.historicalSpeedQuerier

		err := m.ApplyTraffic(c.r, c.liveTraffic, c.historicalSpeed)
		if err != c.expect {
			t.Errorf("Apply traffic on %v (live traffic %v %t, historical speed %v %t), expect %v but got %v", c.r, c.liveTrafficQuerier, c.liveTraffic, c.historicalSpeedQuerier, c.historicalSpeed, c.expect, err)
		}
	}
}

func TestApplyFixedTraffic(t *testing.T) {

	mockFixedSpeed := 100.0
	mockFixedLevel := trafficproxy.TrafficLevel_FREE_FLOW
	mockTraffic := mock.NewFixedTraffic(mockFixedSpeed, mockFixedLevel)

	r := mock.NewOSRMRouteNormal()
	waysCount := len(r.Legs[0].Annotation.Ways)
	appliedLiveTrafficSpeed := make([]float64, waysCount)
	appliedLiveTrafficLevel := make([]int, waysCount)
	appliedBlockIncident := make([]bool, waysCount)
	appliedHistoricalSpeed := make([]float64, waysCount)
	for i := 0; i < waysCount; i++ {
		appliedLiveTrafficSpeed[i] = mockFixedSpeed
		appliedLiveTrafficLevel[i] = int(mockFixedLevel)
		appliedBlockIncident[i] = false
		appliedHistoricalSpeed[i] = mockFixedSpeed
	}

	m, err := appendspeedonly.New(mockTraffic, mockTraffic)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		r                      *route.Route
		liveTraffic            bool
		historicalSpeed        bool
		expectLiveTrafficSpeed []float64
		expectLiveTrafficLevel []int
		expectBlockIncident    []bool
		expectHistoricalSpeed  []float64
	}{
		{mock.NewOSRMRouteNormal(), true, true, appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, appliedHistoricalSpeed},
		{mock.NewOSRMRouteNormal(), false, true, nil, nil, nil, appliedHistoricalSpeed},
		{mock.NewOSRMRouteNormal(), true, false, appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, nil},
		{mock.NewOSRMRouteNormal(), false, false, nil, nil, nil, nil},
	}

	for _, c := range cases {
		if err := m.ApplyTraffic(c.r, c.liveTraffic, c.historicalSpeed); err != nil {
			t.Error(err)
		}

		for _, l := range c.r.Legs {
			if !reflect.DeepEqual(l.Annotation.LiveTrafficSpeed, c.expectLiveTrafficSpeed) {
				t.Errorf("Applied fixed traffic on route %v, expect live traffic speed %v but got %v", c.r, c.expectLiveTrafficSpeed, l.Annotation.LiveTrafficSpeed)
			}
			if !reflect.DeepEqual(l.Annotation.LiveTrafficLevel, c.expectLiveTrafficLevel) {
				t.Errorf("Applied fixed traffic on route %v, expect live traffic level %v but got %v", c.r, c.expectLiveTrafficLevel, l.Annotation.LiveTrafficLevel)
			}
			if !reflect.DeepEqual(l.Annotation.BlockIncident, c.expectBlockIncident) {
				t.Errorf("Applied fixed traffic on route %v, expect live traffic block incident %v but got %v", c.r, c.expectBlockIncident, l.Annotation.BlockIncident)
			}
			if !reflect.DeepEqual(l.Annotation.HistoricalSpeed, c.expectHistoricalSpeed) {
				t.Errorf("Applied fixed traffic on route %v, expect historical speed %v but got %v", c.r, c.expectHistoricalSpeed, l.Annotation.HistoricalSpeed)
			}
		}

	}
}

func TestApplyNormalTraffic(t *testing.T) {

	mockTraffic := mock.NewNormalTraffic()

	appliedLiveTrafficSpeed := []float64{
		float64(float32(6.110000)), // flow uses float32 to store the speed that will lose some precision
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		float64(float32(106.11)),
		float64(float32(106.11)),
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
	}
	appliedLiveTrafficLevel := []int{
		int(trafficproxy.TrafficLevel_SLOW_SPEED), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, int(trafficproxy.TrafficLevel_FREE_FLOW), int(trafficproxy.TrafficLevel_FREE_FLOW), 0, 0,
	}
	appliedBlockIncident := []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false}
	appliedHistoricalSpeed := []float64{
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		20.5,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		trafficapplyingmodel.InvalidTrafficSpeed,
		70.0,
		70.0,
	}

	m, err := appendspeedonly.New(mockTraffic, mockTraffic)
	if err != nil {
		t.Error(err)
	}

	cases := []struct {
		r                      *route.Route
		liveTraffic            bool
		historicalSpeed        bool
		expectLiveTrafficSpeed []float64
		expectLiveTrafficLevel []int
		expectBlockIncident    []bool
		expectHistoricalSpeed  []float64
	}{
		{mock.NewOSRMRouteNormal(), true, true, appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, appliedHistoricalSpeed},
		{mock.NewOSRMRouteNormal(), false, true, nil, nil, nil, appliedHistoricalSpeed},
		{mock.NewOSRMRouteNormal(), true, false, appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, nil},
		{mock.NewOSRMRouteNormal(), false, false, nil, nil, nil, nil},
	}

	for _, c := range cases {
		if err := m.ApplyTraffic(c.r, c.liveTraffic, c.historicalSpeed); err != nil {
			t.Error(err)
		}

		for _, l := range c.r.Legs {
			if !reflect.DeepEqual(l.Annotation.LiveTrafficSpeed, c.expectLiveTrafficSpeed) {
				t.Errorf("Applied traffic on route %v, expect live traffic speed %v but got %v", c.r, c.expectLiveTrafficSpeed, l.Annotation.LiveTrafficSpeed)
			}
			if !reflect.DeepEqual(l.Annotation.LiveTrafficLevel, c.expectLiveTrafficLevel) {
				t.Errorf("Applied traffic on route %v, expect live traffic level %v but got %v", c.r, c.expectLiveTrafficLevel, l.Annotation.LiveTrafficLevel)
			}
			if !reflect.DeepEqual(l.Annotation.BlockIncident, c.expectBlockIncident) {
				t.Errorf("Applied traffic on route %v, expect live traffic block incident %v but got %v", c.r, c.expectBlockIncident, l.Annotation.BlockIncident)
			}
			if !reflect.DeepEqual(l.Annotation.HistoricalSpeed, c.expectHistoricalSpeed) {
				t.Errorf("Applied traffic on route %v, expect historical speed %v but got %v", c.r, c.expectHistoricalSpeed, l.Annotation.HistoricalSpeed)
			}
		}

	}
}
