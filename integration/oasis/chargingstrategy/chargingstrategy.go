package chargingstrategy

import (
	"github.com/Telenav/osrm-backend/integration/util"
	"github.com/golang/glog"
)

type fakeChargingStrategyCreator struct {
	maxEnergyLevel float64
}

// NewFakeChargingStrategyCreator creates fake charging strategy
func NewFakeChargingStrategyCreator(maxEnergyLevel float64) *fakeChargingStrategyCreator {
	return &fakeChargingStrategyCreator{
		maxEnergyLevel: maxEnergyLevel,
	}
}

// @todo:
// - Influence of returning candidate with no charge time and additional energy
// CreateChargingStrategies returns different charging strategy
func (f *fakeChargingStrategyCreator) CreateChargingStrategies() []State {
	return []State{
		State{
			ChargingEnergy: f.maxEnergyLevel * 0.6,
		},
		State{
			ChargingEnergy: f.maxEnergyLevel * 0.8,
		},
		State{
			ChargingEnergy: f.maxEnergyLevel,
		},
	}
}

// Fake charge strategy
// From empty energy:
//                    1 hour charge to 60% of max energy
//                    2 hour charge to 80%, means from 60% ~ 80% need 1 hour
//                    4 hour charge to 100%, means from 80% ~ 100% need 2 hours
func (f *fakeChargingStrategyCreator) EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost {
	sixtyPercentOfMaxEnergy := f.maxEnergyLevel * 0.6
	eightyPercentOfMaxEnergy := f.maxEnergyLevel * 0.8
	noNeedCharge := ChargingCost{
		Duration: 0.0,
	}

	if arrivalEnergy > targetState.ChargingEnergy ||
		util.FloatEquals(targetState.ChargingEnergy, 0.0) {
		return noNeedCharge
	}

	totalTime := 0.0
	currentEnergy := arrivalEnergy
	if arrivalEnergy < sixtyPercentOfMaxEnergy {
		energyNeeded4Stage1 := sixtyPercentOfMaxEnergy - arrivalEnergy
		totalTime += energyNeeded4Stage1 / sixtyPercentOfMaxEnergy * 3600.0
		currentEnergy = sixtyPercentOfMaxEnergy
	}

	if util.FloatEquals(targetState.ChargingEnergy, sixtyPercentOfMaxEnergy) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	if arrivalEnergy < eightyPercentOfMaxEnergy {
		energyNeeded4Stage2 := eightyPercentOfMaxEnergy - currentEnergy
		totalTime += energyNeeded4Stage2 / (eightyPercentOfMaxEnergy - sixtyPercentOfMaxEnergy) * 3600.0
		currentEnergy = eightyPercentOfMaxEnergy
	}
	if util.FloatEquals(targetState.ChargingEnergy, eightyPercentOfMaxEnergy) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	if arrivalEnergy < f.maxEnergyLevel {
		energyNeeded4Stage3 := f.maxEnergyLevel - currentEnergy
		totalTime += energyNeeded4Stage3 / (f.maxEnergyLevel - eightyPercentOfMaxEnergy) * 7200.0
	}

	if util.FloatEquals(targetState.ChargingEnergy, f.maxEnergyLevel) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	glog.Fatalf("Invalid charging state %#v\n", targetState)
	return noNeedCharge
}
