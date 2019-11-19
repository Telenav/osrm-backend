package wayid2nodeids

import (
	"sync"
)

// Mapping handles 'wayID->NodeID,NodeID,NodeID,...' mapping.
type Mapping struct {
	mappingFile   string
	wayID2NodeIDs map[int64][]int64

	ready bool
	mutex sync.RWMutex
}

// NewMappingFrom creates a new Mapping object for 'wayID->NodeID,NodeID,NodeID,...' mapping.
// Currently it only supports mapping file compressed by snappy, e.g. 'wayid2nodeids.csv.snappy'.
func NewMappingFrom(mappingFilePath string) *Mapping {
	m := Mapping{
		mappingFilePath,
		map[int64][]int64{},
		false,
		sync.RWMutex{},
	}
	return &m
}

// Load loads data from file to map in memory, it will returns until the whole load process done.
func (m *Mapping) Load() error {
	defer func() {
		m.mutex.Lock()
		m.ready = true
		m.mutex.Unlock()
	}()
	return m.load()
}

// GetNodeIDs gets nodeIDs mapped by wayID.
func (m *Mapping) GetNodeIDs(wayID int64) []int64 {
	if !m.IsReady() {
		return nil
	}

	nodeIDs, found := m.wayID2NodeIDs[wayID]
	if found {
		return nodeIDs
	}
	return nil
}

// IsReady returns whether the Mapping has been prepared or not.
func (m *Mapping) IsReady() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.ready {
		return true
	}
	return false
}
