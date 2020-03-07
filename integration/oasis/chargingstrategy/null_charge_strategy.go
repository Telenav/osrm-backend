package chargingstrategy

type nullChargeStrategy struct {
}

// NewNullChargeStrategy creates nullChargeStrategy used to bypass unit tests
func NewNullChargeStrategy() *nullChargeStrategy {
	return &nullChargeStrategy{}
}

func (f *nullChargeStrategy) CreateChargingStrategies() []ChargingStatus {
	return []ChargingStatus{}
}

func (f *nullChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState ChargingStatus) ChargingCost {
	return ChargingCost{}
}
