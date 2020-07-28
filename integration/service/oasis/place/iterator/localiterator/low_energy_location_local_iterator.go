package localiterator

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
)

// lowEnergySearchRadius defines search radius for low energy location
const lowEnergySearchRadius = 80000

type lowEnergyLocationIterator struct {
	*basicLocalIterator
}

func newLowEnergyLocationIterator(localFinder place.Finder, location *nav.Location) *lowEnergyLocationIterator {

	obj := &lowEnergyLocationIterator{
		newBasicLocalIterator(localFinder),
	}
	obj.getNearbyChargeStations(nav.Location{
		Lat: location.Lat,
		Lon: location.Lon},
		lowEnergySearchRadius)

	return obj

}
