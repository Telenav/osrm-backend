package localfinder

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

type localStationFinder struct {
	localFinder spatialindexer.Finder
}

// New creates finder based on telenav search web service
func New(localFinder spatialindexer.Finder) *localStationFinder {
	return &localStationFinder{
		localFinder: localFinder,
	}
}

// NewOrigStationFinder creates finder to search for nearby charge stations near orig based on telenav search
func (finder *localStationFinder) NewOrigStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return newOrigStationFinder(finder.localFinder, oasisReq)
}

// NewDestStationFinder creates finder to search for nearby charge stations near destination based on telenav search
func (finder *localStationFinder) NewDestStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return newDestStationFinder(finder.localFinder, oasisReq)
}

// NewLowEnergyLocationStationFinder creates finder to search for nearby charge stations when energy is low based on telenav search
func (finder *localStationFinder) NewLowEnergyLocationStationFinder(location *nav.Location) stationfindertype.NearbyStationsIterator {
	return newLowEnergyLocationLocalFinder(finder.localFinder, location)
}
