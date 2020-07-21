package solution

// StaionSelectionStrategy defines enum of how to select optimal charge stations
type StaionSelectionStrategy int

const (
	// FindChargeStaionsAlongRoute means first calculate a route, then try to find charge stations along the route when energy is low
	FindChargeStaionsAlongRoute = StaionSelectionStrategy(iota) + 1
	// ChargeStaionBasedRouting builds a graph of charge stations and apply shortest-path-algorithm on to it
	ChargeStaionBasedRouting
)
