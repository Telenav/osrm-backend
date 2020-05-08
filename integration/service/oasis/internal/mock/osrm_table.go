package mock

import (
	"github.com/Telenav/osrm-backend/integration/api/osrm"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
)

// 1 * 4
var mockFloatArray1 []float64 = []float64{22.2, 11.1, 33.3, 44.4}
var Mock1To4TableResponse1 = table.Response{
	Durations: [][]float64{
		{mockFloatArray1[0], mockFloatArray1[1], mockFloatArray1[2], mockFloatArray1[3]},
	},
	Distances: [][]float64{
		{mockFloatArray1[0], mockFloatArray1[1], mockFloatArray1[2], mockFloatArray1[3]},
	},
	Sources: []*osrm.Waypoint{
		{
			Name: "orig_location",
		},
	},
	Destinations: []*osrm.Waypoint{
		{
			Name: "station1",
		},
		{
			Name: "station2",
		},
		{
			Name: "station3",
		},
		{
			Name: "station4",
		},
	},
}

// 4 * 1
var mockFloatArray2 []float64 = []float64{66.6, 11.1, 33.3, 33.3}
var Mock4To1TableResponse1 = table.Response{
	Durations: [][]float64{
		{mockFloatArray2[0]},
		{mockFloatArray2[1]},
		{mockFloatArray2[2]},
		{mockFloatArray2[3]},
	},
	Distances: [][]float64{
		{mockFloatArray2[0]},
		{mockFloatArray2[1]},
		{mockFloatArray2[2]},
		{mockFloatArray2[3]},
	},
	Sources: []*osrm.Waypoint{
		{
			Name: "station1",
		},
		{
			Name: "station2",
		},
		{
			Name: "station3",
		},
		{
			Name: "station4",
		},
	},
	Destinations: []*osrm.Waypoint{
		{
			Name: "dest_location",
		},
	},
}

// 4 * 2
var mockFloatArray3 [][]float64 = [][]float64{
	{2, 3},
	{4, 5},
	{6, 7},
	{8, 9},
}

var Mock4To2TableResponse1 = table.Response{
	Durations: [][]float64{
		{mockFloatArray3[0][0], mockFloatArray3[0][1]},
		{mockFloatArray3[1][0], mockFloatArray3[1][1]},
		{mockFloatArray3[2][0], mockFloatArray3[2][1]},
		{mockFloatArray3[3][0], mockFloatArray3[3][1]},
	},
	Distances: [][]float64{
		{mockFloatArray3[0][0], mockFloatArray3[0][1]},
		{mockFloatArray3[1][0], mockFloatArray3[1][1]},
		{mockFloatArray3[2][0], mockFloatArray3[2][1]},
		{mockFloatArray3[3][0], mockFloatArray3[3][1]},
	},
	Sources: []*osrm.Waypoint{
		{
			Name: "station1",
		},
		{
			Name: "station2",
		},
		{
			Name: "station3",
		},
		{
			Name: "station4",
		},
	},
	Destinations: []*osrm.Waypoint{
		{
			Name: "station6",
		},
		{
			Name: "station7",
		},
	},
}

// 2 * 1
var mockFloatArray4 []float64 = []float64{66.6, 11.1}
var Mock2To1TableResponse1 = table.Response{
	Durations: [][]float64{
		{mockFloatArray4[0]},
		{mockFloatArray4[1]},
	},
	Distances: [][]float64{
		{mockFloatArray4[0]},
		{mockFloatArray4[1]},
	},
	Sources: []*osrm.Waypoint{
		{
			Name: "station6",
		},
		{
			Name: "station7",
		},
	},
	Destinations: []*osrm.Waypoint{
		{
			Name: "dest_location",
		},
	},
}
