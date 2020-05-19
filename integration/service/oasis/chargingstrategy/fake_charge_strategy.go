package chargingstrategy

import (
	"math"

	"github.com/Telenav/osrm-backend/integration/util"
)

type fakeChargeStrategy struct {
	maxEnergyLevel                float64
	sixtyPercentOFMaxEnergy       float64
	eightyPercentOfMaxEnergy      float64
	costFrom60PercentTo80Percent  float64
	costFrom60PercentTo100Percent float64
	costFrom80PercentTo100Percent float64
	stateCandidates               []State
}

// NewFakeChargingStrategy creates fake charge strategy
func NewFakeChargingStrategy(maxEnergyLevel float64) *fakeChargeStrategy {
	sixtyPercentOFMaxEnergy := math.Round(maxEnergyLevel*0.6*100) / 100
	eightyPercentOfMaxEnergy := math.Round(maxEnergyLevel*0.8*100) / 100
	maxEnergyLevel = math.Round(maxEnergyLevel*100) / 100
	costFrom60PercentTo80Percent := 3600.0
	costFrom60PercentTo100Percent := 10800.0
	costFrom80PercentTo100Percent := 7200.0
	stateCandidates := []State{
		{
			Energy: sixtyPercentOFMaxEnergy,
		},
		{
			Energy: eightyPercentOfMaxEnergy,
		},
		{
			Energy: maxEnergyLevel,
		},
	}

	return &fakeChargeStrategy{
		maxEnergyLevel:                maxEnergyLevel,
		sixtyPercentOFMaxEnergy:       sixtyPercentOFMaxEnergy,
		eightyPercentOfMaxEnergy:      eightyPercentOfMaxEnergy,
		costFrom60PercentTo80Percent:  costFrom60PercentTo80Percent,
		costFrom60PercentTo100Percent: costFrom60PercentTo100Percent,
		costFrom80PercentTo100Percent: costFrom80PercentTo100Percent,
		stateCandidates:               stateCandidates,
	}
}

// @todo:
// - Influence of returning candidate with no charge time and additional energy
// CreateChargingStates returns different charging strategy
func (f *fakeChargeStrategy) CreateChargingStates() []State {
	return f.stateCandidates
}

var noNeedChargeCost = ChargingCost{
	Duration: 0.0,
}

// Fake charge strategy
// From empty energy:
// charge rule #1:   1 hour charge to 60% of max energy
// charge rule #2:   2 hour charge to 80%, means from 60% ~ 80% need 1 hour
// charge rule #3:   4 hour charge to 100%, means from 80% ~ 100% need 2 hours
func (f *fakeChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost {
	sixtyPercentOfMaxEnergy := f.sixtyPercentOFMaxEnergy
	eightyPercentOfMaxEnergy := f.eightyPercentOfMaxEnergy

	if arrivalEnergy > targetState.Energy ||
		util.Float64Equal(targetState.Energy, 0.0) {
		return noNeedChargeCost
	}

	totalTime := 0.0

	if arrivalEnergy < sixtyPercentOfMaxEnergy {
		energyNeeded4Stage1 := sixtyPercentOfMaxEnergy - arrivalEnergy
		totalTime += energyNeeded4Stage1 / sixtyPercentOfMaxEnergy * 3600.0

		if util.Float64Equal(targetState.Energy, sixtyPercentOfMaxEnergy) {
			return ChargingCost{
				Duration: totalTime,
			}
		} else if util.Float64Equal(targetState.Energy, eightyPercentOfMaxEnergy) {
			return ChargingCost{
				Duration: totalTime + f.costFrom60PercentTo80Percent,
			}
		}
		return ChargingCost{
			Duration: totalTime + f.costFrom60PercentTo100Percent,
		}
	}

	if arrivalEnergy < eightyPercentOfMaxEnergy {
		energyNeeded4Stage2 := eightyPercentOfMaxEnergy - arrivalEnergy
		totalTime += energyNeeded4Stage2 / (eightyPercentOfMaxEnergy - sixtyPercentOfMaxEnergy) * 3600.0

		if util.Float64Equal(targetState.Energy, eightyPercentOfMaxEnergy) {
			return ChargingCost{
				Duration: totalTime,
			}
		}
		return ChargingCost{
			Duration: totalTime + f.costFrom80PercentTo100Percent,
		}
	}

	if arrivalEnergy < f.maxEnergyLevel {
		energyNeeded4Stage3 := f.maxEnergyLevel - arrivalEnergy
		totalTime += energyNeeded4Stage3 / (f.maxEnergyLevel - eightyPercentOfMaxEnergy) * 7200.0

		if util.Float64Equal(targetState.Energy, f.maxEnergyLevel) {
			return ChargingCost{
				Duration: totalTime,
			}
		}
	}

	return noNeedChargeCost
}
