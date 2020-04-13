package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/cloudfinder"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

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
		cloudfinder.LowEnergyLocationCandidateNumber)

	return obj

}

func (localFinder *lowEnergyLocationLocalFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	return localFinder.basicFinder.IterateNearbyStations()
}

func (localFinder *lowEnergyLocationLocalFinder) Stop() {
	localFinder.basicFinder.Stop()
}
