package connectivitymap

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/glog"
)

// Weight represent weight information
type Weight struct {
	Distance float64
	Duration float64
}

// IDAndWeight wraps ID and weight information
type IDAndWeight struct {
	ID     spatialindexer.PointID
	Weight Weight
}

// ID2NearByIDsMap is a mapping between ID and its nearby IDs
type ID2NearByIDsMap map[spatialindexer.PointID][]IDAndWeight

// Connectivity Map used to query connectivity for given placeID
type ConnectivityMap struct {
	id2nearByIDs ID2NearByIDsMap
	maxRange     float64
	statistic    *statistic
}

// New creates ConnectivityMap
func New(maxRange float64) *ConnectivityMap {
	return &ConnectivityMap{
		maxRange:  maxRange,
		statistic: newStatistic(),
	}
}

// Build creates ConnectivityMap
func (cm *ConnectivityMap) Build(iterator spatialindexer.PointsIterator, finder spatialindexer.Finder,
	ranker spatialindexer.Ranker, numOfWorkers int) *ConnectivityMap {
	glog.Info("Start ConnectivityMap's Build().\n")

	cm.id2nearByIDs = newConnectivityMapBuilder(iterator, finder, ranker, cm.maxRange, numOfWorkers).build()
	cm.statistic = cm.statistic.build(cm.id2nearByIDs, cm.maxRange)

	glog.Info("Finished ConnectivityMap's Build().\n")
	return cm
}

// Dump dump ConnectivityMap's content to given folderPath
func (cm *ConnectivityMap) Dump(folderPath string) {
	glog.Info("Start ConnectivityMap's Dump().\n")

	if err := removeAllDumpFiles(folderPath); err != nil {
		glog.Fatalf("removeAllDumpFiles for ConnectivityMap failed with error %+v\n", err)
	}

	if err := serializeConnectivityMap(cm, folderPath); err != nil {
		glog.Fatalf("serializeConnectivityMap failed with error %+v\n", err)
	}

	glog.Infof("Finished ConnectivityMap's Dump() into %s.\n", folderPath)
}

// Load rebuild ConnectivityMap from dumpped data in given folderPath
func (cm *ConnectivityMap) Load(folderPath string) *ConnectivityMap {
	glog.Info("Start ConnectivityMap's Load().\n")

	if err := deSerializeConnectivityMap(cm, folderPath); err != nil {
		glog.Fatalf("deSerializeConnectivityMap failed with error %+v\n", err)
	}

	glog.Infof("Finished ConnectivityMap's Load() from %s.\n", folderPath)
	return cm
}

// QueryConnectivity answers connectivity query for given placeID
// Return true and IDAndWeight array for given placeID, otherwise false and nil
func (cm *ConnectivityMap) QueryConnectivity(placeID spatialindexer.PointID) ([]IDAndWeight, bool) {
	if result, ok := cm.id2nearByIDs[placeID]; ok {
		return result, true
	}
	return nil, false
}

// MaxRange tells the value used to pre-process place data.
// MaxRange means the maximum distance in meters could be reached from current location.
func (cm *ConnectivityMap) MaxRange() float64 {
	return cm.maxRange
}
