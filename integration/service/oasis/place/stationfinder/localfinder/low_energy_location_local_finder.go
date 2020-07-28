package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
)

// lowEnergySearchRadius defines search radius for low energy location
const lowEnergySearchRadius = 80000

type lowEnergyLocationLocalFinder struct {
	*basicLocalFinder
}

func newLowEnergyLocationLocalFinder(localFinder place.Finder, location *nav.Location) *lowEnergyLocationLocalFinder {

	obj := &lowEnergyLocationLocalFinder{
		newBasicLocalFinder(localFinder),
	}
	obj.getNearbyChargeStations(nav.Location{
		Lat: location.Lat,
		Lon: location.Lon},
		lowEnergySearchRadius)

	return obj

}
