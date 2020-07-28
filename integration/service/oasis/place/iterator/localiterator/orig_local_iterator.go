package localiterator

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/golang/glog"
)

type origIterator struct {
	*basicLocalIterator
}

func newOrigIterator(localFinder place.Finder, oasisReq *oasis.Request) *origIterator {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Incorrect oasis request pass into newOrigIterator, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
	}

	obj := &origIterator{
		newBasicLocalIterator(localFinder),
	}
	obj.getNearbyChargeStations(nav.Location{
		Lat: oasisReq.Coordinates[0].Lat,
		Lon: oasisReq.Coordinates[0].Lon},
		oasisReq.CurrRange)

	return obj
}
