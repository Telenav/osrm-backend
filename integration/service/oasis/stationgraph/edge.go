package stationgraph

type neighborer interface {
	neighbors() *[]neighbor
}

type neighbor struct {
	targetNodeID nodeID
	distance     float64
	duration     float64
}
