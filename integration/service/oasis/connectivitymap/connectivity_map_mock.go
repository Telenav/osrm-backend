package connectivitymap

var fakeID2NearByIDsMap1 = ID2NearByIDsMap{
	1: []IDAndDistance{
		IDAndDistance{
			ID:       2,
			Distance: 3,
		},
		IDAndDistance{
			ID:       5,
			Distance: 4,
		},
		IDAndDistance{
			ID:       7,
			Distance: 6,
		},
		IDAndDistance{
			ID:       8,
			Distance: 12,
		},
	},

	2: []IDAndDistance{
		IDAndDistance{
			ID:       1,
			Distance: 3,
		},
		IDAndDistance{
			ID:       7,
			Distance: 23,
		},
	},

	5: []IDAndDistance{
		IDAndDistance{
			ID:       1,
			Distance: 4,
		},
		IDAndDistance{
			ID:       8,
			Distance: 5,
		},
	},

	7: []IDAndDistance{
		IDAndDistance{
			ID:       1,
			Distance: 6,
		},
		IDAndDistance{
			ID:       2,
			Distance: 23,
		},
	},

	8: []IDAndDistance{
		IDAndDistance{
			ID:       5,
			Distance: 5,
		},
		IDAndDistance{
			ID:       1,
			Distance: 12,
		},
	},
}
