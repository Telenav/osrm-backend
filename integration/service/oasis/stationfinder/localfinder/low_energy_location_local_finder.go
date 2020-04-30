package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
)

// lowEnergySearchRadius defines search radius for low energy location
const lowEnergySearchRadius = 80000

type lowEnergyLocationLocalFinder struct {
	*basicLocalFinder
}

func newLowEnergyLocationLocalFinder(localFinder spatialindexer.Finder, location *nav.Location) *lowEnergyLocationLocalFinder {

	obj := &lowEnergyLocationLocalFinder{
		newBasicLocalFinder(localFinder),
	}
	obj.getNearbyChargeStations(spatialindexer.Location{
		Lat: location.Lat,
		Lon: location.Lon},
		lowEnergySearchRadius)

	return obj

}
