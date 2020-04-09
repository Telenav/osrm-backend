package waysnodes

// Writer is the interface that wraps the Write method.
type Writer interface {

	// Write writes wayID->nodeIDs mapping into cache or storage.
	// wayID: is undirected when input, so will always be positive.
	Write(wayID int64, nodeIDs int64) error
}
