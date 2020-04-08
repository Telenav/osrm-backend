package spatialindexer

var mockPlaceInfo1 = []*PointInfo{
	&PointInfo{
		ID: 1,
		Location: Location{
			Lat: 37.355204,
			Lon: -121.953901,
		},
	},
	&PointInfo{
		ID: 2,
		Location: Location{
			Lat: 37.399331,
			Lon: -121.981193,
		},
	},
	&PointInfo{
		ID: 3,
		Location: Location{
			Lat: 37.401948,
			Lon: -121.977384,
		},
	},
	&PointInfo{
		ID: 4,
		Location: Location{
			Lat: 37.407082,
			Lon: -121.991937,
		},
	},
	&PointInfo{
		ID: 5,
		Location: Location{
			Lat: 37.407277,
			Lon: -121.925482,
		},
	},
	&PointInfo{
		ID: 6,
		Location: Location{
			Lat: 37.375024,
			Lon: -121.904706,
		},
	},
	&PointInfo{
		ID: 7,
		Location: Location{
			Lat: 37.359592,
			Lon: -121.914164,
		},
	},
	&PointInfo{
		ID: 8,
		Location: Location{
			Lat: 37.366023,
			Lon: -122.080777,
		},
	},
	&PointInfo{
		ID: 9,
		Location: Location{
			Lat: 37.368453,
			Lon: -122.076400,
		},
	},
	&PointInfo{
		ID: 10,
		Location: Location{
			Lat: 37.373546,
			Lon: -122.068904,
		},
	},
}

// MockPointsIterator implements Finder's interface
type MockFinder struct {
}

// FindNearByPointIDs returns mock result
// It returns 10 places defined in mockPlaceInfo1
func (finder *MockFinder) FindNearByPointIDs(center Location, radius float64, limitCount int) []*PointInfo {
	return mockPlaceInfo1
}

// MockPointsIterator implements PointsIterator's interface
type MockPointsIterator struct {
}

// IteratePoints() iterate places with mock data
func (iterator *MockPointsIterator) IteratePoints() <-chan PointInfo {
	pointInfoC := make(chan PointInfo, len(mockPlaceInfo1))

	go func() {
		for _, item := range mockPlaceInfo1 {
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
func (iterator *MockOneHundredPointsIterator) IteratePoints() <-chan PointInfo {
	pointInfoC := make(chan PointInfo, 100)

	go func() {
		for i := 0; i < 100; i++ {
			id := (PointID)(i + 1000)
			pointInfoC <- PointInfo{
				ID: id,
				Location: Location{
					Lat: 37.398896,
					Lon: -121.976665,
				},
			}
		}

		close(pointInfoC)
	}()

	return pointInfoC
}
