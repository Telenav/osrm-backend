package localiterator

import (
	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
)

type localBasedIteratorGen struct {
	localFinder place.Finder
}

// New creates finder based on telenav search web service
func New(localFinder place.Finder) *localBasedIteratorGen {
	return &localBasedIteratorGen{
		localFinder: localFinder,
	}
}

// NewIterator4Orig creates finder to search for nearby charge stations near orig based on telenav search
func (finder *localBasedIteratorGen) NewIterator4Orig(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator {
	return newOrigIterator(finder.localFinder, oasisReq)
}

// NewIterator4Dest creates finder to search for nearby charge stations near destination based on telenav search
func (finder *localBasedIteratorGen) NewIterator4Dest(oasisReq *oasis.Request) iteratortype.NearbyStationsIterator {
	return newDestIterator(finder.localFinder, oasisReq)
}

// NewIterator4LowEnergyLocation creates finder to search for nearby charge stations when energy is low based on telenav search
func (finder *localBasedIteratorGen) NewIterator4LowEnergyLocation(location *nav.Location) iteratortype.NearbyStationsIterator {
	return newLowEnergyLocationIterator(finder.localFinder, location)
}
