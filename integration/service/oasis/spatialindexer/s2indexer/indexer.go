package s2indexer

import (
	"strconv"
	"time"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/placeloader"
	"github.com/golang/geo/s2"
	"github.com/golang/glog"
)

type cellID2PointIDMap map[s2.CellID][]spatialindexer.PlaceID
type pointID2LocationMap map[spatialindexer.PlaceID]nav.Location

// S2Indexer provide spatial index ability based on google s2
type S2Indexer struct {
	cellID2PointIDs  cellID2PointIDMap
	pointID2Location pointID2LocationMap
}

// NewS2Indexer generates spatial indexer based on google s2
func NewS2Indexer() *S2Indexer {
	return &S2Indexer{
		cellID2PointIDs:  make(cellID2PointIDMap),
		pointID2Location: make(pointID2LocationMap),
	}
}

// Build constructs S2 indexer
func (indexer *S2Indexer) Build(filePath string) *S2Indexer {
	glog.Info("Start Build S2Indexer.\n")
	startTime := time.Now()

	records, err := placeloader.LoadData(filePath)
	if err != nil || len(records) == 0 {
		glog.Error("Failed to Build S2Indexer.\n")
		return nil
	}

	var pointInfos []spatialindexer.PlaceInfo

	for _, record := range records {
		pointInfo := spatialindexer.PlaceInfo{
			ID: elementID2PointID(record.ID),
			Location: nav.Location{
				Lat: record.Lat,
				Lon: record.Lon,
			},
		}
		pointInfos = append(pointInfos, pointInfo)

		indexer.pointID2Location[elementID2PointID(record.ID)] = nav.Location{
			Lat: record.Lat,
			Lon: record.Lon,
		}
	}

	indexer.cellID2PointIDs = build(pointInfos, minS2Level, maxS2Level)

	glog.Info("Finished Build() S2Indexer.\n")
	glog.Infof("Build S2Indexer takes %f seconds.\n", time.Since(startTime).Seconds())
	return indexer
}

// Load S2Indexer's data from contents recorded in folder
func (indexer *S2Indexer) Load(folderPath string) *S2Indexer {
	glog.Info("Start S2Indexer's Load().\n")

	if err := deSerializeS2Indexer(indexer, folderPath); err != nil {
		glog.Errorf("Load S2Indexer's data from folder %s failed, err=%v\n", folderPath, err)
		return nil
	}

	glog.Infof("Finished S2Indexer's Load() from folder %s.\n", folderPath)
	return indexer
}

// Dump S2Indexer's content into folderPath
func (indexer *S2Indexer) Dump(folderPath string) {
	glog.Info("Start S2Indexer's Dump().\n")

	if err := serializeS2Indexer(indexer, folderPath); err != nil {
		glog.Errorf("Dump S2Indexer's data to folder %s failed, err=%v\n", folderPath, err)
	}

	glog.Infof("Finished S2Indexer's Dump() to folder %s.\n", folderPath)
}

// IteratePlaces returns PlaceInfo in channel
// It implements interface of PlacesIterator
func (indexer *S2Indexer) IteratePlaces() <-chan spatialindexer.PlaceInfo {
	placesC := make(chan spatialindexer.PlaceInfo, len(indexer.pointID2Location))
	go func() {
		for pointID, location := range indexer.pointID2Location {
			placesC <- spatialindexer.PlaceInfo{
				ID:       pointID,
				Location: location,
			}
		}
		close(placesC)
	}()

	return placesC
}

// FindNearByPlaceIDs returns nearby places for given center and conditions
func (indexer *S2Indexer) FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*spatialindexer.PlaceInfo {
	if !indexer.isInitialized() {
		glog.Warning("S2Indexer is empty, try to Build() with correct input file first.\n")
		return nil
	}

	results := queryNearByPlaces(indexer, center, radius)
	if limitCount != spatialindexer.UnlimitedCount && len(results) > limitCount {
		results = results[:limitCount]
	}

	return results
}

// GetLocation returns *nav.Location for given placeID
// Returns nil if given placeID is not found
func (indexer *S2Indexer) GetLocation(placeID string) *nav.Location {
	id, err := strconv.Atoi(placeID)
	if err != nil {
		glog.Errorf("Incorrect station ID passed to NearByStationQuery %+v, got error %#v", placeID, err)
		return nil
	}
	if location, ok := indexer.pointID2Location[(spatialindexer.PlaceID)(id)]; ok {
		return &nav.Location{
			Lat: location.Lat,
			Lon: location.Lon,
		}
	}

	return nil
}

//TODO codebear801 This function should be replaced by GetLocation
func (indexer S2Indexer) getPointLocationByPointID(id spatialindexer.PlaceID) (nav.Location, bool) {
	location, ok := indexer.pointID2Location[id]
	return location, ok
}

func (indexer S2Indexer) getPointIDsByS2CellID(cellid s2.CellID) ([]spatialindexer.PlaceID, bool) {
	pointIDs, ok := indexer.cellID2PointIDs[cellid]
	return pointIDs, ok
}

func (indexer S2Indexer) isInitialized() bool {
	return len(indexer.cellID2PointIDs) != 0 && len(indexer.pointID2Location) != 0
}

func elementID2PointID(id int64) spatialindexer.PlaceID {
	return (spatialindexer.PlaceID)(id)
}
