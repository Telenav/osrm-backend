package chargingstrategy

type fakeChargingStrategyCreator struct {
	arrivalEnergyLevel float64
	maxEnergyLevel     float64
}

// NewFakeChargingStrategyCreator creates fake charging strategy
func NewFakeChargingStrategyCreator(arrivalEnergyLevel, maxEnergyLevel float64) *fakeChargingStrategyCreator {
	return &fakeChargingStrategyCreator{
		arrivalEnergyLevel: arrivalEnergyLevel,
		maxEnergyLevel:     maxEnergyLevel,
	}
}

// @todo:
// - Influence of returning candidate with no charge time and additional energy

// CreateChargingStrategies returns different charging strategy
// Initial implementation:
// From empty energy:
//                    1 hour charge to 60% of max energy
//                    2 hour charge to 80%, means from 60% ~ 80% need 1 hour
//                    4 hour charge to 100%, means from 80% ~ 100% need 2 hours
func (f *fakeChargingStrategyCreator) CreateChargingStrategies() []ChargingStrategy {
	sixtyPercentOfMaxEnergy := f.maxEnergyLevel * 0.6
	eightyPercentOfMaxEnergy := f.maxEnergyLevel * 0.8
	var result []ChargingStrategy

	if f.arrivalEnergyLevel < sixtyPercentOfMaxEnergy {
		energy4Stage1 := sixtyPercentOfMaxEnergy - f.arrivalEnergyLevel
		time4Stage1 := energy4Stage1 / sixtyPercentOfMaxEnergy * 3600.0
		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage1,
			ChargingEnergy: energy4Stage1,
		})

		energy4Stage2 := eightyPercentOfMaxEnergy - sixtyPercentOfMaxEnergy + energy4Stage1
		time4Stage2 := 3600.0 + time4Stage1
		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage2,
			ChargingEnergy: energy4Stage2,
		})

		energy4Stage3 := f.maxEnergyLevel - sixtyPercentOfMaxEnergy + energy4Stage1
		time4Stage3 := 7200.0 + 3600.0 + time4Stage1
		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage3,
			ChargingEnergy: energy4Stage3,
		})
	} else if f.arrivalEnergyLevel < eightyPercentOfMaxEnergy {
		energy4Stage2 := eightyPercentOfMaxEnergy - f.arrivalEnergyLevel
		time4Stage2 := energy4Stage2 / (eightyPercentOfMaxEnergy - sixtyPercentOfMaxEnergy) * 3600

		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage2,
			ChargingEnergy: energy4Stage2,
		})

		energy4Stage3 := f.maxEnergyLevel - eightyPercentOfMaxEnergy + energy4Stage2
		time4Stage3 := 7200.0 + time4Stage2
		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage3,
			ChargingEnergy: energy4Stage3,
		})
	} else {
		energy4Stage3 := f.maxEnergyLevel - f.arrivalEnergyLevel
		time4Stage3 := energy4Stage3 / (f.maxEnergyLevel - eightyPercentOfMaxEnergy) * 7200
		result = append(result, ChargingStrategy{
			ChargingTime:   time4Stage3,
			ChargingEnergy: energy4Stage3,
		})
	}

	return result
}
