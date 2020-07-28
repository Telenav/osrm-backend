package chargingstrategy

import (
	"reflect"
	"testing"
)

func TestSimpleChargingStrategyCreator(t *testing.T) {
	cases := []struct {
		arrivalEnergyLevel float64
		maxEnergyLevel     float64
		expectResult       []ChargingCost
	}{
		{
			10000,
			50000,
			[]ChargingCost{
				ChargingCost{
					Duration: 2400.0,
				},
				ChargingCost{
					Duration: 6000.0,
				},
				ChargingCost{
					Duration: 13200.0,
				},
			},
		},
		{
			32000,
			50000,
			[]ChargingCost{
				ChargingCost{
					Duration: 0.0,
				},
				ChargingCost{
					Duration: 2880.0,
				},
				ChargingCost{
					Duration: 10080.0,
				},
			},
		},
		{
			41000,
			50000,
			[]ChargingCost{
				ChargingCost{
					Duration: 0.0,
				},
				ChargingCost{
					Duration: 0.0,
				},
				ChargingCost{
					Duration: 6480.0,
				},
			},
		},
	}

	for _, c := range cases {
		var actualResult []ChargingCost
		strategy := NewSimpleChargingStrategy(c.maxEnergyLevel)
		for _, state := range strategy.CreateChargingStates() {
			actualResult = append(actualResult, strategy.EvaluateCost(c.arrivalEnergyLevel, state))
		}
		if !reflect.DeepEqual(actualResult, c.expectResult) {
			t.Errorf("parse case %#v, expect\n %#v but got\n %#v", c, c.expectResult, actualResult)
		}
	}
}
