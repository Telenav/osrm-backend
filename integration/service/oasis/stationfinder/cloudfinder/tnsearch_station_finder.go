package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

type tnSearchStationFinder struct {
	sc *searchconnector.TNSearchConnector
}

// New creates finder based on telenav search web service
func New(sc *searchconnector.TNSearchConnector) *tnSearchStationFinder {
	return &tnSearchStationFinder{
		sc: sc,
	}
}

// NewOrigStationFinder creates finder to search for nearby charge stations near orig based on telenav search
func (finder *tnSearchStationFinder) NewOrigStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return NewOrigStationFinder(finder.sc, oasisReq)
}

// NewDestStationFinder creates finder to search for nearby charge stations near destination based on telenav search
func (finder *tnSearchStationFinder) NewDestStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return NewDestStationFinder(finder.sc, oasisReq)
}

// NewLowEnergyLocationStationFinder creates finder to search for nearby charge stations when energy is low based on telenav search
func (finder *tnSearchStationFinder) NewLowEnergyLocationStationFinder(location *nav.Location) stationfindertype.NearbyStationsIterator {
	return NewLowEnergyLocationStationFinder(finder.sc, location)
}
