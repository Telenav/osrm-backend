// Package modelfactory provides support to create specified traffic applying model applier.
package modelfactory

import (
	"fmt"
	"strings"

	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/appendonly"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/preferlivetraffic"
	"github.com/Telenav/osrm-backend/integration/traffic"
	"github.com/golang/glog"
)

// AvailableModelNames returns name list of all available traffic applying models.
func AvailableModelNames() []string {

	n := []string{}
	for k := range availableModels {
		n = append(n, k)
	}
	return n
}

// DefaultModelName returns the name of default traffic applying model.
func DefaultModelName() string {
	return defaultModel
}

// NewApplier creates a new specified traffic applier.
func NewApplier(name string, l traffic.LiveTrafficQuerier, h traffic.HistoricalSpeedQuerier) (trafficapplyingmodel.Applier, error) {
	if _, ok := availableModels[name]; !ok {
		return nil, fmt.Errorf("invalid traffic applying model %s, options: %s", name, strings.Join(AvailableModelNames(), ","))
	}

	var applier trafficapplyingmodel.Applier
	var err error
	switch name {
	case appendonly.Name():
		applier, err = appendonly.New(l, h)
	case preferlivetraffic.Name():
		applier, err = preferlivetraffic.New(l, h)
	default:
		glog.Fatalf("missed traffic applying model: %s", name)
	}

	if err != nil {
		return nil, err
	}
	return applier, nil
}

var (
	availableModels = map[string]struct{}{
		appendonly.Name():        struct{}{},
		preferlivetraffic.Name(): struct{}{},
	}
	defaultModel = preferlivetraffic.Name()
)
