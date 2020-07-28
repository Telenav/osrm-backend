package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/golang/glog"
)

type origStationLocalFinder struct {
	*basicLocalFinder
}

func newOrigStationFinder(localFinder place.Finder, oasisReq *oasis.Request) *origStationLocalFinder {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Incorrect oasis request pass into newOrigStationFinder, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
	}

	obj := &origStationLocalFinder{
		newBasicLocalFinder(localFinder),
	}
	obj.getNearbyChargeStations(nav.Location{
		Lat: oasisReq.Coordinates[0].Lat,
		Lon: oasisReq.Coordinates[0].Lon},
		oasisReq.CurrRange)

	return obj
}
