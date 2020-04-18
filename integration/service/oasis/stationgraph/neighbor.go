package stationgraph

// Neighborer defines adjacent list for given node
type Neighborer interface {

	// Neighbors returns neighbor information for given node
	Neighbors() []*neighbor
}

type neighbor struct {
	targetNodeID nodeID
	distance     float64
	duration     float64
}
