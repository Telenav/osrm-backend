package mock

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

// MockFinder implements Finder's interface
type MockFinder struct {
}

// FindNearByPlaceIDs returns mock result
// It returns 10 places defined in MockPlaceInfo1
func (finder *MockFinder) FindNearByPlaceIDs(center nav.Location, radius float64, limitCount int) []*entity.PlaceInfo {
	return MockPlaceInfo1
}

// MockPlaceInfo1 contains 10 PlaceInfo items
var MockPlaceInfo1 = []*entity.PlaceInfo{
	{
		ID: 1,
		Location: &nav.Location{
			Lat: 37.355204,
			Lon: -121.953901,
		},
	},
	{
		ID: 2,
		Location: &nav.Location{
			Lat: 37.399331,
			Lon: -121.981193,
		},
	},
	{
		ID: 3,
		Location: &nav.Location{
			Lat: 37.401948,
			Lon: -121.977384,
		},
	},
	{
		ID: 4,
		Location: &nav.Location{
			Lat: 37.407082,
			Lon: -121.991937,
		},
	},
	{
		ID: 5,
		Location: &nav.Location{
			Lat: 37.407277,
			Lon: -121.925482,
		},
	},
	{
		ID: 6,
		Location: &nav.Location{
			Lat: 37.375024,
			Lon: -121.904706,
		},
	},
	{
		ID: 7,
		Location: &nav.Location{
			Lat: 37.359592,
			Lon: -121.914164,
		},
	},
	{
		ID: 8,
		Location: &nav.Location{
			Lat: 37.366023,
			Lon: -122.080777,
		},
	},
	{
		ID: 9,
		Location: &nav.Location{
			Lat: 37.368453,
			Lon: -122.076400,
		},
	},
	{
		ID: 10,
		Location: &nav.Location{
			Lat: 37.373546,
			Lon: -122.068904,
		},
	},
}
