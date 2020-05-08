package mock

import "github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"

// MockFinder implements Finder's interface
type MockFinder struct {
}

// FindNearByPointIDs returns mock result
// It returns 10 places defined in MockPlaceInfo1
func (finder *MockFinder) FindNearByPointIDs(center spatialindexer.Location, radius float64, limitCount int) []*spatialindexer.PointInfo {
	return MockPlaceInfo1
}

// MockPlaceInfo1 contains 10 PointInfo items
var MockPlaceInfo1 = []*spatialindexer.PointInfo{
	{
		ID: 1,
		Location: spatialindexer.Location{
			Lat: 37.355204,
			Lon: -121.953901,
		},
	},
	{
		ID: 2,
		Location: spatialindexer.Location{
			Lat: 37.399331,
			Lon: -121.981193,
		},
	},
	{
		ID: 3,
		Location: spatialindexer.Location{
			Lat: 37.401948,
			Lon: -121.977384,
		},
	},
	{
		ID: 4,
		Location: spatialindexer.Location{
			Lat: 37.407082,
			Lon: -121.991937,
		},
	},
	{
		ID: 5,
		Location: spatialindexer.Location{
			Lat: 37.407277,
			Lon: -121.925482,
		},
	},
	{
		ID: 6,
		Location: spatialindexer.Location{
			Lat: 37.375024,
			Lon: -121.904706,
		},
	},
	{
		ID: 7,
		Location: spatialindexer.Location{
			Lat: 37.359592,
			Lon: -121.914164,
		},
	},
	{
		ID: 8,
		Location: spatialindexer.Location{
			Lat: 37.366023,
			Lon: -122.080777,
		},
	},
	{
		ID: 9,
		Location: spatialindexer.Location{
			Lat: 37.368453,
			Lon: -122.076400,
		},
	},
	{
		ID: 10,
		Location: spatialindexer.Location{
			Lat: 37.373546,
			Lon: -122.068904,
		},
	},
}
