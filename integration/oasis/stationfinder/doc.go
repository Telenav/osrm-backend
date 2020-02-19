// Package stationfinder provide functionality to find nearby charge stations and
// related algorithm.
// Finders:
// - origStationFinder holds logic for how to find reachable charge stations
//   based on current energy level.
// - destStationFinder holds logic for how to find reachable charge stations
//   based on safe energy level and distance to nearest charge station(todo).
// - lowEnergyLocationStationFinder holds logic for how to find reachable
//   charge station near certain location.

// Algorithm:
// - Each finder provide iterator to iterate charge station candidates.
// - The choice of channel as response makes algorithm could be asynchronous func.
// - FindOverlapBetweenStations provide functionality to find overlap
//   between two iterator.
// - CalcCostBetweenChargeStationsPair provide functionality to calculate
//   cost between two group of charge stations which could construct a new
//   graph as edges.
package stationfinder
