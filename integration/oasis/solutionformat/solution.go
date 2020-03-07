package solutionformat

// Solution contains summary and selected charge stations
type Solution struct {
	Distance       float64
	Duration       float64
	RemainingRage  float64
	Weight         float64
	ChargeStations []*ChargeStation
}

// ChargeStation contains all information related with specific charge station
type ChargeStation struct {
	Location      Location
	StationID     string
	ArrivalEnergy float64
	WaitTime      float64
	ChargeTime    float64
	ChargeRange   float64
}

// Location defines the geo location of a station
type Location struct {
	Lat float64
	Lon float64
}
