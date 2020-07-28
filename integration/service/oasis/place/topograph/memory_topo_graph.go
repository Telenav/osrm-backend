package topograph

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/golang/glog"
)

// ID2NearByIDsMap is a mapping between ID and its nearby IDs
type ID2NearByIDsMap map[entity.PlaceID][]*entity.TransferInfo

// Connectivity Map used to query connectivity for given placeID
type MemoryTopoGraph struct {
	id2nearByIDs ID2NearByIDsMap
	maxRange     float64
	statistic    *statistic
}

// New creates MemoryTopoGraph
func New(maxRange float64) *MemoryTopoGraph {
	return &MemoryTopoGraph{
		maxRange:  maxRange,
		statistic: newStatistic(),
	}
}

// Build creates MemoryTopoGraph
func (cm *MemoryTopoGraph) Build(iterator place.Iterator, finder place.Finder,
	ranker place.Ranker, numOfWorkers int) *MemoryTopoGraph {
	glog.Info("Start MemoryTopoGraph's Build().\n")

	cm.id2nearByIDs = newConnectivityMapBuilder(iterator, finder, ranker, cm.maxRange, numOfWorkers).build()
	cm.statistic = cm.statistic.build(cm.id2nearByIDs, cm.maxRange)

	glog.Info("Finished MemoryTopoGraph's Build().\n")
	return cm
}

// Dump dump MemoryTopoGraph's content to given folderPath
func (cm *MemoryTopoGraph) Dump(folderPath string) {
	glog.Info("Start MemoryTopoGraph's Dump().\n")

	if err := removeAllDumpFiles(folderPath); err != nil {
		glog.Fatalf("removeAllDumpFiles for MemoryTopoGraph failed with error %+v\n", err)
	}

	if err := serializeConnectivityMap(cm, folderPath); err != nil {
		glog.Fatalf("serializeConnectivityMap failed with error %+v\n", err)
	}

	glog.Infof("Finished MemoryTopoGraph's Dump() into %s.\n", folderPath)
}

// Load rebuild MemoryTopoGraph from dumpped data in given folderPath
func (cm *MemoryTopoGraph) Load(folderPath string) *MemoryTopoGraph {
	glog.Info("Start MemoryTopoGraph's Load().\n")

	if err := deSerializeConnectivityMap(cm, folderPath); err != nil {
		glog.Fatalf("deSerializeConnectivityMap failed with error %+v\n", err)
	}

	glog.Infof("Finished MemoryTopoGraph's Load() from %s.\n", folderPath)
	return cm
}

// QueryConnectivity answers connectivity query for given placeID
// Return true and IDAndWeight array for given placeID, otherwise false and nil
func (cm *MemoryTopoGraph) QueryConnectivity(placeID entity.PlaceID) ([]*entity.TransferInfo, bool) {
	if result, ok := cm.id2nearByIDs[placeID]; ok {
		return result, true
	}
	return nil, false
}

// MaxRange tells the value used to pre-process place data.
// MaxRange means the maximum distance in meters could be reached from current location.
func (cm *MemoryTopoGraph) MaxRange() float64 {
	return cm.maxRange
}
