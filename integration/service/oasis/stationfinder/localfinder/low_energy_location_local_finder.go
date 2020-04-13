package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

// lowEnergySearchRadius defines search radius for low energy location
const lowEnergySearchRadius = 80000

type lowEnergyLocationLocalFinder struct {
	basicFinder *basicLocalFinder
}

func newLowEnergyLocationLocalFinder(localFinder spatialindexer.Finder, location *nav.Location) *lowEnergyLocationLocalFinder {

	obj := &lowEnergyLocationLocalFinder{
		basicFinder: newBasicLocalFinder(localFinder),
	}
	obj.basicFinder.getNearbyChargeStations(spatialindexer.Location{
		Lat: location.Lat,
		Lon: location.Lon},
		lowEnergySearchRadius)

	return obj

}

// // NearbyStationsIterator provides channel which contains near by station information for low energy point
func (localFinder *lowEnergyLocationLocalFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	return localFinder.basicFinder.IterateNearbyStations()
}

// Stop stops functionality of finder
func (localFinder *lowEnergyLocationLocalFinder) Stop() {
	localFinder.basicFinder.Stop()
}
