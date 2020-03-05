package chargingstrategy

type fakeChargingStrategyCreator struct {
	maxEnergyLevel float64
}

// NewFakeChargingStrategyCreator creates fake charging strategy
func NewFakeChargingStrategyCreator(maxEnergyLevel float64) *fakeChargingStrategyCreator {
	return &fakeChargingStrategyCreator{
		maxEnergyLevel: maxEnergyLevel,
	}
}

// CreateChargingStrategies returns different charging strategy
// Initial implementation: 1 hour charge for 60% of max energy,
//                         2 hour charge for 80%
//                         4 hour charge for 100%
func (f *fakeChargingStrategyCreator) CreateChargingStrategies() []ChargingStrategy {

	return []ChargingStrategy{
		ChargingStrategy{
			ChargingTime:   3600,
			ChargingEnergy: f.maxEnergyLevel * 0.6,
		},
		ChargingStrategy{
			ChargingTime:   7200,
			ChargingEnergy: f.maxEnergyLevel * 0.8,
		},
		ChargingStrategy{
			ChargingTime:   14400,
			ChargingEnergy: f.maxEnergyLevel,
		},
	}
}
