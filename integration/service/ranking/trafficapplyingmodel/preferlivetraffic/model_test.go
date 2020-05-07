package preferlivetraffic_test

import (
	"reflect"
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/preferlivetraffic"
	"github.com/Telenav/osrm-backend/integration/traffic"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
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
		{mock.NewOSRMRouteNoDataSourceName(), mock.EmptyTraffic{}, mock.EmptyTraffic{}, true, true, trafficapplyingmodel.ErrEmptyAnnotationMetadata},
	}

	m, err := preferlivetraffic.New(mock.EmptyTraffic{}, mock.EmptyTraffic{})
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

func TestApplyNormalTraffic(t *testing.T) {
	mockTraffic := mock.NewNormalTrafficNoBlock()

	r := mock.NewOSRMRouteNormal()
	appliedDataSourceNamesBoth := append(r.Legs[0].Annotation.Metadata.DataSourceNames, trafficapplyingmodel.SourceNameHistoricalSpeed, trafficapplyingmodel.SourceNameLiveTraffic)
	appliedDataSourceNamesHistoricalSpeed := append(r.Legs[0].Annotation.Metadata.DataSourceNames, trafficapplyingmodel.SourceNameHistoricalSpeed)
	appliedDataSourceNamesLiveTraffic := append(r.Legs[0].Annotation.Metadata.DataSourceNames, trafficapplyingmodel.SourceNameLiveTraffic)
	appliedDataSourceNamesNone := append(r.Legs[0].Annotation.Metadata.DataSourceNames)

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
	appliedBlockIncident := []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false}
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

	liveTrafficDurationWeightChanges := (15.936241/6.110000 - 15.936241/8.9) + (44.431011/106.11 - 44.431011/8.9) + (69.128936/106.11 - 69.128936/9)
	historicalSpeedDurationWeightChanges := (40.00651/20.5 - 40.00651/8.9) + (57.603303/70.0 - 57.603303/9) + (11.706858/70.0 - 11.706858/9)

	m, err := preferlivetraffic.New(mockTraffic, mockTraffic)
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
		expectSpeed            []float64
		expectDataSourceNames  []string
		expectDataSources      []int
		expectDuration         float64
		expectWeight           float64
	}{
		{
			mock.NewOSRMRouteNormal(),
			true, true,
			appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, appliedHistoricalSpeed,
			[]float64{float64(float32(6.110000)), 8.9, 20.5, 8.9, 14.3, 14.2, 14.4, 14.6, 14.7, 14, 14.4, 14.7, 14, 14.6, float64(float32(106.11)), float64(float32(106.11)), 70.0, 70.0},
			appliedDataSourceNamesBoth,
			[]int{2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 1},
			60.3 + liveTrafficDurationWeightChanges + historicalSpeedDurationWeightChanges,
			60.3 + liveTrafficDurationWeightChanges + historicalSpeedDurationWeightChanges,
		},
		{
			mock.NewOSRMRouteNormal(),
			false, true,
			nil, nil, nil, appliedHistoricalSpeed,
			[]float64{8.9, 8.9, 20.5, 8.9, 14.3, 14.2, 14.4, 14.6, 14.7, 14, 14.4, 14.7, 14, 14.6, 8.9, 9, 70.0, 70.0},
			appliedDataSourceNamesHistoricalSpeed,
			[]int{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
			60.3 + historicalSpeedDurationWeightChanges,
			60.3 + historicalSpeedDurationWeightChanges,
		},
		{
			mock.NewOSRMRouteNormal(),
			true, false,
			appliedLiveTrafficSpeed, appliedLiveTrafficLevel, appliedBlockIncident, nil,
			[]float64{float64(float32(6.110000)), 8.9, 8.9, 8.9, 14.3, 14.2, 14.4, 14.6, 14.7, 14, 14.4, 14.7, 14, 14.6, float64(float32(106.11)), float64(float32(106.11)), 9, 9},
			appliedDataSourceNamesLiveTraffic,
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
			60.3 + liveTrafficDurationWeightChanges,
			60.3 + liveTrafficDurationWeightChanges,
		},
		{
			mock.NewOSRMRouteNormal(),
			false, false,
			nil, nil, nil, nil,
			[]float64{8.9, 8.9, 8.9, 8.9, 14.3, 14.2, 14.4, 14.6, 14.7, 14, 14.4, 14.7, 14, 14.6, 8.9, 9, 9, 9},
			appliedDataSourceNamesNone,
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			60.3,
			60.3,
		},
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
			if !reflect.DeepEqual(l.Annotation.Speed, c.expectSpeed) {
				t.Errorf("Applied traffic on route %v, expect speed %v but got %v", c.r, c.expectSpeed, l.Annotation.Speed)
			}
			if !reflect.DeepEqual(l.Annotation.Metadata.DataSourceNames, c.expectDataSourceNames) {
				t.Errorf("Applied traffic on route %v, expect data source names %v but got %v", c.r, c.expectDataSourceNames, l.Annotation.Metadata.DataSourceNames)
			}
			if !reflect.DeepEqual(l.Annotation.DataSources, c.expectDataSources) {
				t.Errorf("Applied traffic on route %v, expect data sources %v but got %v", c.r, c.expectDataSources, l.Annotation.DataSources)
			}
			if !roughFloatEqual(l.Weight, c.expectWeight) || !roughFloatEqual(l.Duration, c.expectDuration) {
				t.Errorf("Applied traffic on route %v, expect weight, duration %f,%f but got %f,%f", c.r, c.expectWeight, c.expectDuration, l.Weight, l.Duration)
			}

		}

		if !roughFloatEqual(c.r.Weight, c.expectWeight) || !roughFloatEqual(c.r.Duration, c.expectDuration) {
			t.Errorf("Applied traffic on route %v, expect weight, duration %f,%f but got %f,%f", c.r, c.expectWeight, c.expectDuration, c.r.Weight, c.r.Duration)
		}

	}
}

func roughFloatEqual(a, b float64) bool {
	percision := 0.1 // only 0.1 percision
	if (a-b) < percision && (b-a) < percision {
		return true
	}
	return false
}
