package chargingstrategy

type fakeChargeStrategy struct {
}

// NewFakeChargeStrategy creates fakeChargeStrategy used to bypass unit tests
func NewFakeChargeStrategy() *fakeChargeStrategy {
	return &fakeChargeStrategy{}
}

func (f *fakeChargeStrategy) CreateChargingStates() []State {
	return []State{}
}

func (f *fakeChargeStrategy) EvaluateCost(arrivalEnergy float64, targetState State) ChargingCost {
	return ChargingCost{}
}
