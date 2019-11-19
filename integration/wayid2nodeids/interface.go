package wayid2nodeids

import (
	"github.com/Telenav/osrm-backend/integration/nodebasededge"
)

// Mapping handles 'wayID->NodeID,NodeID,NodeID,...' mapping.
type Mapping struct {
	mappingFile   string
	wayID2NodeIDs map[int64][]int64
	edge2WayID    map[nodebasededge.Edge]int64
	nodeIDs       map[int64]struct{}
}

// NewMappingFrom creates a new Mapping object for 'wayID->NodeID,NodeID,NodeID,...' mapping.
// Currently it only supports mapping file compressed by snappy, e.g. 'wayid2nodeids.csv.snappy'.
func NewMappingFrom(mappingFilePath string) Mapping {
	m := Mapping{
		mappingFilePath, map[int64][]int64{}, map[nodebasededge.Edge]int64{}, map[int64]struct{}{},
	}
	return m
}

// Load loads data from file to map in memory, it will returns until the whole load process done.
func (m *Mapping) Load() error {
	return m.load()
}

// GetNodeIDs gets nodeIDs mapped by wayID.
func (m Mapping) GetNodeIDs(wayID int64) []int64 {
	nodeIDs, found := m.wayID2NodeIDs[wayID]
	if found {
		return nodeIDs
	}
	return nil
}

// GetWayID returns wayID corresponding to Edge.
// The second return bool indicates whether it's found or not.
func (m Mapping) GetWayID(edge nodebasededge.Edge) (int64, bool) {

	if wayID, found := m.edge2WayID[edge]; found {
		return wayID, true
	}

	if wayID, found := m.edge2WayID[edge.Reverse()]; found {
		return -wayID, true
	}

	return 0, false
}
