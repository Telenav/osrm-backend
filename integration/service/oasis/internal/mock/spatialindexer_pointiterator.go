package mock

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

// MockPlacesIterator implements PlacesIterator's interface
type MockPlacesIterator struct {
}

// IteratePlaces() iterate places with mock data
func (iterator *MockPlacesIterator) IteratePlaces() <-chan spatialindexer.PlaceInfo {
	pointInfoC := make(chan spatialindexer.PlaceInfo, len(MockPlaceInfo1))

	go func() {
		for _, item := range MockPlaceInfo1 {
			pointInfoC <- *item
		}

		close(pointInfoC)
	}()

	return pointInfoC
}

// MockOneHundredPlacesIterator implements PlacesIterator's interface
type MockOneHundredPlacesIterator struct {
}

// IteratePlaces() iterate places with mock data.
// It returns {ID:1000, fixed location}, {ID:1001, fixed location}, ... {ID:1099, fixed location}
func (iterator *MockOneHundredPlacesIterator) IteratePlaces() <-chan spatialindexer.PlaceInfo {
	pointInfoC := make(chan spatialindexer.PlaceInfo, 100)

	go func() {
		for i := 0; i < 100; i++ {
			id := (spatialindexer.PlaceID)(i + 1000)
			pointInfoC <- spatialindexer.PlaceInfo{
				ID: id,
				Location: nav.Location{
					Lat: 37.398896,
					Lon: -121.976665,
				},
			}
		}

		close(pointInfoC)
	}()

	return pointInfoC
}
