package connectivitymap

var fakeID2NearByIDsMap1 = ID2NearByIDsMap{
	1: []IDAndWeight{
		{
			ID: 2,
			Weight: Weight{
				Distance: 3,
				Duration: 3,
			},
		},
		{
			ID: 5,
			Weight: Weight{
				Distance: 4,
				Duration: 4,
			},
		},
		{
			ID: 7,
			Weight: Weight{
				Distance: 6,
				Duration: 61,
			},
		},
		{
			ID: 8,
			Weight: Weight{
				Distance: 12,
				Duration: 12,
			},
		},
	},

	2: []IDAndWeight{
		{
			ID: 1,
			Weight: Weight{
				Distance: 3,
				Duration: 3,
			},
		},
		{
			ID: 7,
			Weight: Weight{
				Distance: 23,
				Duration: 23,
			},
		},
	},

	5: []IDAndWeight{
		{
			ID: 1,
			Weight: Weight{
				Distance: 4,
				Duration: 4,
			},
		},
		{
			ID: 8,
			Weight: Weight{
				Distance: 5,
				Duration: 5,
			},
		},
	},

	7: []IDAndWeight{
		{
			ID: 1,
			Weight: Weight{
				Distance: 6,
				Duration: 6,
			},
		},
		{
			ID: 2,
			Weight: Weight{
				Distance: 23,
				Duration: 23,
			},
		},
	},

	8: []IDAndWeight{
		{
			ID: 5,
			Weight: Weight{
				Distance: 5,
				Duration: 5,
			},
		},
		{
			ID: 1,
			Weight: Weight{
				Distance: 12,
				Duration: 12,
			},
		},
	},
}

var fakeID2NearByIDsMap2 = ID2NearByIDsMap{
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

// MockConnectivityMap constructs simple connectivity map for integration test
var MockConnectivityMap = ConnectivityMap{
	id2nearByIDs: fakeID2NearByIDsMap2,
}
