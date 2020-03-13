package chargingstrategy

import (
	"github.com/Telenav/osrm-backend/integration/util"
	"github.com/golang/glog"
)

type fakeChargeStrategy struct {
	maxEnergyLevel float64
}

// NewFakeChargingStrategy creates fake charge strategy
func NewFakeChargingStrategy(maxEnergyLevel float64) *fakeChargeStrategy {
	return &fakeChargeStrategy{
		maxEnergyLevel: maxEnergyLevel,
	}
}

// @todo:
// - Influence of returning candidate with no charge time and additional energy
// CreateChargingStates returns different charging strategy
func (f *fakeChargeStrategy) CreateChargingStates() []State {
	return []State{
		State{
			Energy: f.maxEnergyLevel * 0.6,
		},
		State{
			Energy: f.maxEnergyLevel * 0.8,
		},
		State{
			Energy: f.maxEnergyLevel,
		},
	}
}

// Fake charge strategy
// From empty energy:
//                    1 hour charge to 60% of max energy
//                    2 hour charge to 80%, means from 60% ~ 80% need 1 hour
//                    4 hour charge to 100%, means from 80% ~ 100% need 2 hours
func (f *fakeChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost {
	sixtyPercentOfMaxEnergy := f.maxEnergyLevel * 0.6
	eightyPercentOfMaxEnergy := f.maxEnergyLevel * 0.8
	noNeedCharge := ChargingCost{
		Duration: 0.0,
	}

	if arrivalEnergy > targetState.Energy ||
		util.FloatEquals(targetState.Energy, 0.0) {
		return noNeedCharge
	}

	totalTime := 0.0
	currentEnergy := arrivalEnergy
	if arrivalEnergy < sixtyPercentOfMaxEnergy {
		energyNeeded4Stage1 := sixtyPercentOfMaxEnergy - arrivalEnergy
		totalTime += energyNeeded4Stage1 / sixtyPercentOfMaxEnergy * 3600.0
		currentEnergy = sixtyPercentOfMaxEnergy
	}

	if util.FloatEquals(targetState.Energy, sixtyPercentOfMaxEnergy) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	if arrivalEnergy < eightyPercentOfMaxEnergy {
		energyNeeded4Stage2 := eightyPercentOfMaxEnergy - currentEnergy
		totalTime += energyNeeded4Stage2 / (eightyPercentOfMaxEnergy - sixtyPercentOfMaxEnergy) * 3600.0
		currentEnergy = eightyPercentOfMaxEnergy
	}
	if util.FloatEquals(targetState.Energy, eightyPercentOfMaxEnergy) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	if arrivalEnergy < f.maxEnergyLevel {
		energyNeeded4Stage3 := f.maxEnergyLevel - currentEnergy
		totalTime += energyNeeded4Stage3 / (f.maxEnergyLevel - eightyPercentOfMaxEnergy) * 7200.0
	}

	if util.FloatEquals(targetState.Energy, f.maxEnergyLevel) {
		return ChargingCost{
			Duration: totalTime,
		}
	}

	glog.Fatalf("Invalid charging state %#v\n", targetState)
	return noNeedCharge
}
