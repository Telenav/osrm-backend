package stationgraph

type stationGraph struct {
	g             *graph
	stationID2Int map[string]uint32
}

// @todo
// func NewStationGraph(c chan stationfinder.WeightBetweenNeighbors) *stationGraph {
// 	return nil
// }
