package s2indexer

import (
	"github.com/Telenav/osrm-backend/integration/service/spatialindexer"
	"github.com/golang/geo/s2"
)

// https://s2geometry.io/resources/s2cell_statistics.html
// Level = 6 means average area size is 20754.64km2
// Level = 9 means average area size is 324.29km2
const minS2Level = 9

// Level = 19 means average area size is 309.27km2
// Level = 20 means average area size is 77.32km2
const maxS2Level = 20

func build(points []spatialindexer.PointInfo, minLevel, maxLevel int) map[s2.CellID][]spatialindexer.PointID {
	pointID2CellIDs := make(map[spatialindexer.PointID][]s2.CellID)
	cellID2PointIDs := make(map[s2.CellID][]spatialindexer.PointID)

	for _, p := range points {
		leafCellID := s2.CellFromLatLng(s2.LatLngFromDegrees(p.Location.Lat, p.Location.Lon)).ID()

		var cellIDs []s2.CellID
		// For level = 30, its parent equal to current
		// So no need append leafCellID into cellIDs outside of for loop
		for i := leafCellID.Level(); i >= minLevel; i-- {
			if i > maxLevel {
				continue
			}

			parentCellID := leafCellID.Parent(i)
			cellIDs = append(cellIDs, parentCellID)
		}

		pointID2CellIDs[p.ID] = cellIDs
	}

	for pointID, cellIDs := range pointID2CellIDs {
		for _, cellID := range cellIDs {
			if _, ok := cellID2PointIDs[cellID]; !ok {
				var pointIDs []spatialindexer.PointID
				cellID2PointIDs[cellID] = pointIDs
			}

			cellID2PointIDs[cellID] = append(cellID2PointIDs[cellID], pointID)
		}
	}

	return cellID2PointIDs
}
