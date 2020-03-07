package chargingstrategy

import (
	"reflect"
	"testing"
)

func TestFakeChargingStrategyCreator(t *testing.T) {
	cases := []struct {
		arrivalEnergyLevel float64
		maxEnergyLevel     float64
		expectResult       []ChargingStrategy
	}{
		{
			10000,
			50000,
			[]ChargingStrategy{
				ChargingStrategy{
					ChargingTime:   2400,
					ChargingEnergy: 20000,
				},
				ChargingStrategy{
					ChargingTime:   6000,
					ChargingEnergy: 30000,
				},
				ChargingStrategy{
					ChargingTime:   13200,
					ChargingEnergy: 40000,
				},
			},
		},
		{
			32000,
			50000,
			[]ChargingStrategy{
				ChargingStrategy{
					ChargingTime:   2880,
					ChargingEnergy: 8000,
				},
				ChargingStrategy{
					ChargingTime:   10080,
					ChargingEnergy: 18000,
				},
			},
		},
		{
			41000,
			50000,
			[]ChargingStrategy{
				ChargingStrategy{
					ChargingTime:   6480,
					ChargingEnergy: 9000,
				},
			},
		},
	}

	for _, c := range cases {
		actualResult := NewFakeChargingStrategyCreator(c.arrivalEnergyLevel, c.maxEnergyLevel).CreateChargingStrategies()
		if !reflect.DeepEqual(actualResult, c.expectResult) {
			t.Errorf("parse case %#v, expect %#v but got %#v", c, c.expectResult, actualResult)
		}
	}
}
