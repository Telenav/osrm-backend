package chargingstrategy

// ChargingStrategy contains charging realted information
type ChargingStrategy struct {
	ChargingTime   float64
	ChargingEnergy float64
}

// ChargingStrategyCreator defines interface related with creation of charging strategy
type ChargingStrategyCreator interface {

	// CreateChargingStrategies creates charge strategies which could be used by other algorithm
	CreateChargingStrategies() []ChargingStrategy
}
