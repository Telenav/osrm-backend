package mock

import "github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"

// MockPointsIterator implements PointsIterator's interface
type MockPointsIterator struct {
}

// IteratePoints() iterate places with mock data
func (iterator *MockPointsIterator) IteratePoints() <-chan spatialindexer.PointInfo {
	pointInfoC := make(chan spatialindexer.PointInfo, len(MockPlaceInfo1))

	go func() {
		for _, item := range MockPlaceInfo1 {
			pointInfoC <- *item
		}

		close(pointInfoC)
	}()

	return pointInfoC
}

// MockOneHundredPointsIterator implements PointsIterator's interface
type MockOneHundredPointsIterator struct {
}

// IteratePoints() iterate places with mock data.
// It returns {ID:1000, fixed location}, {ID:1001, fixed location}, ... {ID:1099, fixed location}
func (iterator *MockOneHundredPointsIterator) IteratePoints() <-chan spatialindexer.PointInfo {
	pointInfoC := make(chan spatialindexer.PointInfo, 100)

	go func() {
		for i := 0; i < 100; i++ {
			id := (spatialindexer.PointID)(i + 1000)
			pointInfoC <- spatialindexer.PointInfo{
				ID: id,
				Location: spatialindexer.Location{
					Lat: 37.398896,
					Lon: -121.976665,
				},
			}
		}

		close(pointInfoC)
	}()

	return pointInfoC
}
