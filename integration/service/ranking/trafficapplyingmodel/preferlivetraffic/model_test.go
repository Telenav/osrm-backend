package preferlivetraffic_test

import (
	"testing"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/internal/mock"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/preferlivetraffic"
	"github.com/Telenav/osrm-backend/integration/traffic"
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
