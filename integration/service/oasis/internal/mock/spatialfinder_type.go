package mock

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

// MockChargeStationInfo1 mocks array of *ChargeStationInfo which is compatible with MockPlaceInfo1
var MockChargeStationInfo1 = []*stationfindertype.ChargeStationInfo{
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[0].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[0].Location.Lat,
			Lon: MockPlaceInfo1[0].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[1].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[1].Location.Lat,
			Lon: MockPlaceInfo1[1].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[2].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[2].Location.Lat,
			Lon: MockPlaceInfo1[2].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[3].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[3].Location.Lat,
			Lon: MockPlaceInfo1[3].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[4].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[4].Location.Lat,
			Lon: MockPlaceInfo1[4].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[5].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[5].Location.Lat,
			Lon: MockPlaceInfo1[5].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[6].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[6].Location.Lat,
			Lon: MockPlaceInfo1[6].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[7].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[7].Location.Lat,
			Lon: MockPlaceInfo1[7].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[8].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[8].Location.Lat,
			Lon: MockPlaceInfo1[8].Location.Lon,
		},
		Err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(MockPlaceInfo1[9].ID), 10),
		Location: nav.Location{
			Lat: MockPlaceInfo1[9].Location.Lat,
			Lon: MockPlaceInfo1[9].Location.Lon,
		},
		Err: nil,
	},
}

var NeighborInfoArray0 = []stationfindertype.NeighborInfo{
	{
		FromID: stationfindertype.OrigLocationIDStr,
		FromLocation: nav.Location{
			Lat: 1.1,
			Lon: 1.1,
		},
		ToID: "station1",
		ToLocation: nav.Location{
			Lat: 32.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 22.2,
			Distance: 22.2,
		},
	},
	{
		FromID: stationfindertype.OrigLocationIDStr,
		FromLocation: nav.Location{
			Lat: 1.1,
			Lon: 1.1,
		},
		ToID: "station2",
		ToLocation: nav.Location{
			Lat: -32.333,
			Lon: -122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 11.1,
			Distance: 11.1,
		},
	},
	{
		FromID: stationfindertype.OrigLocationIDStr,
		FromLocation: nav.Location{
			Lat: 1.1,
			Lon: 1.1,
		},
		ToID: "station3",
		ToLocation: nav.Location{
			Lat: 32.333,
			Lon: -122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 33.3,
			Distance: 33.3,
		},
	},
	{
		FromID: stationfindertype.OrigLocationIDStr,
		FromLocation: nav.Location{
			Lat: 1.1,
			Lon: 1.1,
		},
		ToID: "station4",
		ToLocation: nav.Location{
			Lat: -32.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 44.4,
			Distance: 44.4,
		},
	},
}

var NeighborInfoArray1 = []stationfindertype.NeighborInfo{
	{
		FromID: "station1",
		FromLocation: nav.Location{
			Lat: 32.333,
			Lon: 122.333,
		},
		ToID: "station6",
		ToLocation: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 2,
			Distance: 2,
		},
	},
	{
		FromID: "station1",
		FromLocation: nav.Location{
			Lat: 32.333,
			Lon: 122.333,
		},
		ToID: "station7",
		ToLocation: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 3,
			Distance: 3,
		},
	},
	{
		FromID: "station2",
		FromLocation: nav.Location{
			Lat: -32.333,
			Lon: -122.333,
		},
		ToID: "station6",
		ToLocation: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 4,
			Distance: 4,
		},
	},
	{
		FromID: "station2",
		FromLocation: nav.Location{
			Lat: -32.333,
			Lon: -122.333,
		},
		ToID: "station7",
		ToLocation: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 5,
			Distance: 5,
		},
	},
	{
		FromID: "station3",
		FromLocation: nav.Location{
			Lat: 32.333,
			Lon: -122.333,
		},
		ToID: "station6",
		ToLocation: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 6,
			Distance: 6,
		},
	},
	{
		FromID: "station3",
		FromLocation: nav.Location{
			Lat: 32.333,
			Lon: -122.333,
		},
		ToID: "station7",
		ToLocation: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 7,
			Distance: 7,
		},
	},
	{
		FromID: "station4",
		FromLocation: nav.Location{
			Lat: -32.333,
			Lon: 122.333,
		},
		ToID: "station6",
		ToLocation: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 8,
			Distance: 8,
		},
	},
	{
		FromID: "station4",
		FromLocation: nav.Location{
			Lat: -32.333,
			Lon: 122.333,
		},
		ToID: "station7",
		ToLocation: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
		Weight: stationfindertype.Weight{
			Duration: 9,
			Distance: 9,
		},
	},
}

var NeighborInfoArray2 = []stationfindertype.NeighborInfo{
	{
		FromID: "station6",
		FromLocation: nav.Location{
			Lat: 30.333,
			Lon: 122.333,
		},
		ToID: stationfindertype.DestLocationIDStr,
		ToLocation: nav.Location{
			Lat: 4.4,
			Lon: 4.4,
		},
		Weight: stationfindertype.Weight{
			Duration: 66.6,
			Distance: 66.6,
		},
	},
	{
		FromID: "station7",
		FromLocation: nav.Location{
			Lat: -10.333,
			Lon: 122.333,
		},
		ToID: stationfindertype.DestLocationIDStr,
		ToLocation: nav.Location{
			Lat: 4.4,
			Lon: 4.4,
		},
		Weight: stationfindertype.Weight{
			Duration: 11.1,
			Distance: 11.1,
		},
	},
}
