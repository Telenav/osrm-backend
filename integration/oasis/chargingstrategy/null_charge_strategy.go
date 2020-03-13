package chargingstrategy

type nullChargeStrategy struct {
}

// NewNullChargeStrategy creates nullChargeStrategy used to bypass unit tests
func NewNullChargeStrategy() *nullChargeStrategy {
	return &nullChargeStrategy{}
}

func (f *nullChargeStrategy) CreateChargingStates() []State {
	return []State{}
}

func (f *nullChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost {
	return ChargingCost{}
}
