package connectivitymap

import "github.com/Telenav/osrm-backend/integration/service/oasis/internal/common"

var fakeID2NearByIDsMap3 = ID2NearByIDsMap{
	1: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 2,
			},
			Weight: &common.Weight{
				Distance: 1,
				Duration: 1,
			},
		},
	},
	2: []*common.RankedPlaceInfo{
		{
			PlaceInfo: common.PlaceInfo{
				ID: 3,
			},
			Weight: &common.Weight{
				Distance: 2,
				Duration: 2,
			},
		},
	},
}

// todo codebear801 removes MockConnectivityMap which expose internal implementation
// MockConnectivityMap constructs simple connectivity map for integration test
var MockConnectivityMap = ConnectivityMap{
	id2nearByIDs: fakeID2NearByIDsMap3,
}
