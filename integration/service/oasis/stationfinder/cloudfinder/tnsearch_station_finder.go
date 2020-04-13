package cloudfinder

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/nav"
	"github.com/Telenav/osrm-backend/integration/pkg/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/searchconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/stationfinder/stationfindertype"
)

type cloudStationFinder struct {
	sc *searchconnector.TNSearchConnector
}

// New creates finder based on telenav search web service
func New(sc *searchconnector.TNSearchConnector) *cloudStationFinder {
	return &cloudStationFinder{
		sc: sc,
	}
}

// NewOrigStationFinder creates finder to search for nearby charge stations near orig based on telenav search
func (finder *cloudStationFinder) NewOrigStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return NewOrigStationFinder(finder.sc, oasisReq)
}

// NewDestStationFinder creates finder to search for nearby charge stations near destination based on telenav search
func (finder *cloudStationFinder) NewDestStationFinder(oasisReq *oasis.Request) stationfindertype.NearbyStationsIterator {
	return NewDestStationFinder(finder.sc, oasisReq)
}

// NewLowEnergyLocationStationFinder creates finder to search for nearby charge stations when energy is low based on telenav search
func (finder *cloudStationFinder) NewLowEnergyLocationStationFinder(location *nav.Location) stationfindertype.NearbyStationsIterator {
	return NewLowEnergyLocationStationFinder(finder.sc, location)
}
