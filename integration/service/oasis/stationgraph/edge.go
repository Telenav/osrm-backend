package stationgraph

type edgeID struct {
	fromNodeID nodeID
	toNodeID   nodeID
}

type edgeMetric struct {
	distance float64
	duration float64
}

type edge struct {
	edgeId     edgeID
	edgeMetric *edgeMetric
}
