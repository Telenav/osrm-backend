package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

type origStationLocalFinder struct {
	basicFinder *basicLocalFinder
}

func newOrigStationFinder(localFinder spatialindexer.Finder, oasisReq *oasis.Request) *origStationLocalFinder {
	if len(oasisReq.Coordinates) != 2 {
		glog.Errorf("Incorrect oasis request pass into newOrigStationFinder, len(oasisReq.Coordinates) should be 2 but got %d.\n", len(oasisReq.Coordinates))
	}

	obj := &origStationLocalFinder{
		basicFinder: newBasicLocalFinder(localFinder),
	}
	obj.basicFinder.getNearbyChargeStations(spatialindexer.Location{
		Lat: oasisReq.Coordinates[0].Lat,
		Lon: oasisReq.Coordinates[0].Lon},
		oasisReq.CurrRange)

	return obj
}

func (localFinder *origStationLocalFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	return localFinder.basicFinder.IterateNearbyStations()
}

func (localFinder *origStationLocalFinder) Stop() {
	localFinder.basicFinder.Stop()
}
