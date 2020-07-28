package s2indexer

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/golang/glog"
)

const maxCellCount = 200

const s2EarthRadiusInMeters = 6371010.0

func queryNearByS2Cells(point nav.Location, radiusInMeters float64) []s2.CellID {
	regionCover := &s2.RegionCoverer{
		MinLevel: minS2Level,
		MaxLevel: maxS2Level,
		MaxCells: maxCellCount}
	center := s2.PointFromLatLng(s2.LatLngFromDegrees(point.Lat, point.Lon))
	radius := (s1.Angle)(radiusInMeters / s2EarthRadiusInMeters)
	region := s2.Region(s2.CapFromCenterAngle(center, radius))
	cellUnion := regionCover.Covering(region)

	return ([]s2.CellID)(cellUnion)
}

func queryNearByPlaces(indexer *S2Indexer, point nav.Location, radius float64) []*entity.PlaceWithLocation {
	var result []*entity.PlaceWithLocation

	cellIDs := queryNearByS2Cells(point, radius)

	for _, cellID := range cellIDs {
		pointIDs, hasCellID := indexer.getPointIDsByS2CellID(cellID)
		if !hasCellID {
			continue
		}

		for _, pointID := range pointIDs {
			location, hasPointID := indexer.getPointLocationByPointID(pointID)
			if !hasPointID {
				glog.Errorf("In queryNearByPlaces, use incorrect pointID %v to query S2Indexer\n", pointID)
				continue
			}

			result = append(result, &entity.PlaceWithLocation{
				ID:       pointID,
				Location: &location,
			})
		}
	}

	return result
}

func generateDebugInfo4Query(point nav.Location, radius float64, cellIDs []s2.CellID) {
	glog.Infof("During spatial_query, point = %+v, radius = %v {", point, radius)
	for _, cellID := range cellIDs {
		glog.Infof("%s,", cellID.ToToken())
	}
	glog.Info("}\n")
}

func generateDebugInfo4CellIDs(cellIDs []s2.CellID) {
	glog.Info("=================================\n")
	glog.Info("generateDebugInfo4CellIDs\n")

	for _, cellID := range cellIDs {
		glog.Infof("CellID value = %d(uint64), string = %v, token = %s, level = %d\n",
			(uint64)(cellID), cellID, cellID.ToToken(), cellID.Level())
	}

	glog.Info("=================================\n")
}

// Generate sample url like: http://s2.sidewalklabs.com/regioncoverer/?cells=89c2584b54,89c2584d,89c25852c,89c259a5,89c259ac,89c259b4,89c259bc,89c259c7,89c259c9,89c259ca4
func generateDebugURL(cellIDs []s2.CellID) string {
	var url string

	if len(cellIDs) == 0 {
		return url
	}

	url += "http://s2.sidewalklabs.com/regioncoverer/?cells="
	for _, cellID := range cellIDs {
		url += cellID.ToToken() + ","
	}

	return url
}
