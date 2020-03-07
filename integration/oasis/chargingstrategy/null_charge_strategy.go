package chargingstrategy

type nullChargeStrategy struct {
}

// NewNullChargeStrategy creates nullChargeStrategy used to bypass unit tests
func NewNullChargeStrategy() *nullChargeStrategy {
	return &nullChargeStrategy{}
}

func (f *nullChargeStrategy) CreateChargingStrategies() []ChargingStrategy {
	return []ChargingStrategy{}
}

func (f *nullChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState ChargingStrategy) ChargingCost {
	return ChargingCost{}
}
