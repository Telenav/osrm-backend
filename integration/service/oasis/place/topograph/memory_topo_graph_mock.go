package topograph

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
)

var fakeID2NearByIDsMap3 = ID2NearByIDsMap{
	1: []*entity.TransferInfo{
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 2,
				Location: &nav.Location{
					Lat: 37.399331,
					Lon: -121.981193,
				},
			},
			Weight: &entity.Weight{
				Distance: 1,
				Duration: 1,
			},
		},
	},
	2: []*entity.TransferInfo{
		{
			PlaceWithLocation: entity.PlaceWithLocation{
				ID: 3,
				Location: &nav.Location{
					Lat: 37.401948,
					Lon: -121.977384,
				},
			},
			Weight: &entity.Weight{
				Distance: 2,
				Duration: 2,
			},
		},
	},
}

// todo codebear801 removes MockConnectivityMap which expose internal implementation
// MockConnectivityMap constructs simple connectivity map for integration test
var MockConnectivityMap = MemoryTopoGraph{
	id2nearByIDs: fakeID2NearByIDsMap3,
}
