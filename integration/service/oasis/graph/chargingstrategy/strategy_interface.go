package chargingstrategy

// State contains charging related information
type State struct {
	Energy float64
}

// ChargingCost represents the cost needed to reach certain states
type ChargingCost struct {
	Duration float64
	// Later could add money usage, etc
}

// Creator creates charge states based on different charge strategy
type Creator interface {
	// CreateChargingStates creates charge States which could be used by other algorithm
	CreateChargingStates() []State
}

// Evaluator calculate cost from given status(energy, etc) to target State
type Evaluator interface {
	// EvaluateCost accepts current status and target status and returns cost needed
	EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost
}

// Strategy defines interface related with creation and evaluation
type Strategy interface {
	Creator
	Evaluator
}
