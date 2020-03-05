package chargingstrategy

// ChargingStrategy contains charging realted information
type ChargingStrategy struct {
	ChargingTime   float64
	ChargingEnergy float64
}

type ChargingStrategyCreator interface {
	CreateChargingStrategies() []ChargingStrategy
}
