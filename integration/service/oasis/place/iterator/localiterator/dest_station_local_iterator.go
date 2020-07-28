package localiterator

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/golang/glog"
)

type destIterator struct {
	*basicLocalIterator
}

func newDestIterator(localFinder place.Finder, oasisReq *oasis.Request) *destIterator {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Try to create newOrigIterator use incorrect oasis request, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
		return nil
	}
	if oasisReq.MaxRange <= oasisReq.SafeLevel {
		glog.Errorf("Try to create newOrigIterator use incorrect oasis request, SafeLevel should be smaller than MaxRange.\n")
		return nil
	}

	obj := &destIterator{
		newBasicLocalIterator(localFinder),
	}
	obj.getNearbyChargeStations(nav.Location{
		Lat: oasisReq.Coordinates[1].Lat,
		Lon: oasisReq.Coordinates[1].Lon},
		oasisReq.MaxRange-oasisReq.SafeLevel)

	return obj
}
