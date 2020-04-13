package stationfindertype

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

// MockChargeStationInfo1 mocks array of *ChargeStationInfo which is compatible with spatialindexer.MockPlaceInfo1
var MockChargeStationInfo1 = []*ChargeStationInfo{
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[0].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[0].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[0].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[1].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[1].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[1].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[2].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[2].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[2].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[3].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[3].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[3].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[4].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[4].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[4].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[5].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[5].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[5].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[6].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[6].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[6].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[7].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[7].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[7].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[8].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[8].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[8].Location.Lon,
		},
		err: nil,
	},
	{
		ID: strconv.FormatInt((int64)(spatialindexer.MockPlaceInfo1[9].ID), 10),
		Location: nav.Location{
			Lat: spatialindexer.MockPlaceInfo1[9].Location.Lat,
			Lon: spatialindexer.MockPlaceInfo1[9].Location.Lon,
		},
		err: nil,
	},
}
