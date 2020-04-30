package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/golang/glog"
)

type destStationLocalFinder struct {
	*basicLocalFinder
}

func newDestStationFinder(localFinder spatialindexer.Finder, oasisReq *oasis.Request) *destStationLocalFinder {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Try to create newOrigStationFinder use incorrect oasis request, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
		return nil
	}
	if oasisReq.MaxRange <= oasisReq.SafeLevel {
		glog.Errorf("Try to create newOrigStationFinder use incorrect oasis request, SafeLevel should be smaller than MaxRange.\n")
		return nil
	}

	obj := &destStationLocalFinder{
		newBasicLocalFinder(localFinder),
	}
	obj.getNearbyChargeStations(spatialindexer.Location{
		Lat: oasisReq.Coordinates[1].Lat,
		Lon: oasisReq.Coordinates[1].Lon},
		oasisReq.MaxRange-oasisReq.SafeLevel)

	return obj
}
