package table

import (
	"github.com/Telenav/osrm-backend/integration/api/osrm/osrmtype"
)

// 1 * 4
var mockFloatArray1 []float64 = []float64{22.2, 11.1, 33.3, 44.4}
var Mock1ToNTableResponse1 Response = Response{
	Durations: [][]float64{
		[]float64{mockFloatArray1[0], mockFloatArray1[1], mockFloatArray1[2], mockFloatArray1[3]},
	},
	Distances: [][]float64{
		[]float64{mockFloatArray1[0], mockFloatArray1[1], mockFloatArray1[2], mockFloatArray1[3]},
	},
	Sources: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "orig_location",
		},
	},
	Destinations: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "station1",
		},
		&osrmtype.Waypoint{
			Name: "station2",
		},
		&osrmtype.Waypoint{
			Name: "station3",
		},
		&osrmtype.Waypoint{
			Name: "station4",
		},
	},
}

// 4 * 1
var mockFloatArray2 []float64 = []float64{66.6, 11.1, 33.3, 33.3}
var MockNTo1TableResponse1 Response = Response{
	Durations: [][]float64{
		[]float64{mockFloatArray2[0]},
		[]float64{mockFloatArray2[1]},
		[]float64{mockFloatArray2[2]},
		[]float64{mockFloatArray2[3]},
	},
	Distances: [][]float64{
		[]float64{mockFloatArray2[0]},
		[]float64{mockFloatArray2[1]},
		[]float64{mockFloatArray2[2]},
		[]float64{mockFloatArray2[3]},
	},
	Sources: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "station1",
		},
		&osrmtype.Waypoint{
			Name: "station2",
		},
		&osrmtype.Waypoint{
			Name: "station3",
		},
		&osrmtype.Waypoint{
			Name: "station4",
		},
	},
	Destinations: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "dest_location",
		},
	},
}

// 4 * 2
var mockFloatArray3 [][]float64 = [][]float64{
	[]float64{2, 3},
	[]float64{4, 5},
	[]float64{6, 7},
	[]float64{8, 9},
}

var Mock4To2TableResponse1 Response = Response{
	Durations: [][]float64{
		[]float64{mockFloatArray3[0][0], mockFloatArray3[0][1]},
		[]float64{mockFloatArray3[1][0], mockFloatArray3[1][1]},
		[]float64{mockFloatArray3[2][0], mockFloatArray3[2][1]},
		[]float64{mockFloatArray3[3][0], mockFloatArray3[3][1]},
	},
	Distances: [][]float64{
		[]float64{mockFloatArray3[0][0], mockFloatArray3[0][1]},
		[]float64{mockFloatArray3[1][0], mockFloatArray3[1][1]},
		[]float64{mockFloatArray3[2][0], mockFloatArray3[2][1]},
		[]float64{mockFloatArray3[3][0], mockFloatArray3[3][1]},
	},
	Sources: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "station1",
		},
		&osrmtype.Waypoint{
			Name: "station2",
		},
		&osrmtype.Waypoint{
			Name: "station3",
		},
		&osrmtype.Waypoint{
			Name: "station4",
		},
	},
	Destinations: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "station6",
		},
		&osrmtype.Waypoint{
			Name: "station7",
		},
	},
}

// 2 * 1
var mockFloatArray4 []float64 = []float64{66.6, 11.1}
var Mock2To1TableResponse1 Response = Response{
	Durations: [][]float64{
		[]float64{mockFloatArray4[0]},
		[]float64{mockFloatArray4[1]},
	},
	Distances: [][]float64{
		[]float64{mockFloatArray4[0]},
		[]float64{mockFloatArray4[1]},
	},
	Sources: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "station6",
		},
		&osrmtype.Waypoint{
			Name: "station7",
		},
	},
	Destinations: []*osrmtype.Waypoint{
		&osrmtype.Waypoint{
			Name: "dest_location",
		},
	},
}
