package clouditerator

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/Telenav/osrm-backend/integration/util/searchconnector"
)

type cloudBasedIteratorGen struct {
	sc *searchconnector.TNSearchConnector
}

// New creates IteratorGenerator based on telenav search web service
func New(sc *searchconnector.TNSearchConnector) *cloudBasedIteratorGen {
	return &cloudBasedIteratorGen{
		sc: sc,
	}
}

// NewIterator4Orig creates iterator to search for nearby charge stations near orig based on telenav search
func (c *cloudBasedIteratorGen) NewIterator4Orig(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator {
	return newOrigIterator(c.sc, oasisReq)
}

// NewIterator4Dest creates iterator to search for nearby charge stations near destination based on telenav search
func (c *cloudBasedIteratorGen) NewIterator4Dest(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator {
	return newDestIterator(c.sc, oasisReq)
}

// NewIterator4LowEnergyLocation creates iterator to search for nearby charge stations when energy is low based on telenav search
func (c *cloudBasedIteratorGen) NewIterator4LowEnergyLocation(location *nav.Location) iteratortype.NearbyStationsIterator {
	return newLowEnergyLocationIterator(c.sc, location)
}
