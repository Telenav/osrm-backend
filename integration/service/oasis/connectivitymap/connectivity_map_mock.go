package connectivitymap

var fakeID2NearByIDsMap3 = ID2NearByIDsMap{
	1: []IDAndWeight{
		{
			ID: 2,
			Weight: Weight{
				Distance: 1,
				Duration: 1,
			},
		},
	},
	2: []IDAndWeight{
		{
			ID: 3,
			Weight: Weight{
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
