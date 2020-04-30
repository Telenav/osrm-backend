package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/glog"
)

type origStationLocalFinder struct {
	*basicLocalFinder
}

func newOrigStationFinder(localFinder spatialindexer.Finder, oasisReq *oasis.Request) *origStationLocalFinder {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Incorrect oasis request pass into newOrigStationFinder, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
	}

	obj := &origStationLocalFinder{
		newBasicLocalFinder(localFinder),
	}
	obj.getNearbyChargeStations(spatialindexer.Location{
		Lat: oasisReq.Coordinates[0].Lat,
		Lon: oasisReq.Coordinates[0].Lon},
		oasisReq.CurrRange)

	return obj
}
