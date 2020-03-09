package chargingstrategy

// State contains charging related information
type State struct {
	ChargingEnergy float64
}

// ChargingCost represents the cost needed to reach certain states
type ChargingCost struct {
	Duration float64
	// Later could add money usage, etc
}

// ChargingStrategyCreator defines interface related with creation of charging strategy
type ChargingStrategyCreator interface {

	// CreateChargingStrategies creates charge strategies which could be used by other algorithm
	CreateChargingStrategies() []State

	// EvaluateCost accepts current status and target status and returns cost needed
	EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost
}
