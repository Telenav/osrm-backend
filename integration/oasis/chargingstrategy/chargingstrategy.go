package chargingstrategy

type ChargingStrategy struct {
	ChargingTime   float64
	ChargingEnergy float64
}

// CreateChargingStrategies returns different charging strategy
// Initial implementation: 1 hour charge for 60% of max energy,
//                         2 hour charge for 80%
//                         4 hour charge for 100%
func CreateChargingStrategies(maxEnergyLevel float64) []ChargingStrategy {

	return []ChargingStrategy{
		ChargingStrategy{
			ChargingTime:   3600,
			ChargingEnergy: maxEnergyLevel * 0.6,
		},
		ChargingStrategy{
			ChargingTime:   7200,
			ChargingEnergy: maxEnergyLevel * 0.8,
		},
		ChargingStrategy{
			ChargingTime:   14400,
			ChargingEnergy: maxEnergyLevel,
		},
	}
}
