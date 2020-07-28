package localiterator

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/iterator/iteratortype"
	"github.com/golang/glog"
)

const defaultIteratorCount = 5
const defaultChargeStaionChannelSize = 500

type basicLocalIterator struct {
	localFinder place.Finder
	placesInfo  []*entity.PlaceWithLocation
	requests    chan chan *iteratortype.ChargeStationInfo
	stop        chan bool
}

func newBasicLocalIterator(localFinder place.Finder) *basicLocalIterator {
	bf := &basicLocalIterator{
		localFinder: localFinder,
		requests:    make(chan chan *iteratortype.ChargeStationInfo, defaultIteratorCount),
		stop:        make(chan bool),
	}
	go bf.serveRequest()
	return bf
}

func (bf *basicLocalIterator) getNearbyChargeStations(center nav.Location, radius float64) {
	bf.placesInfo = bf.localFinder.FindNearByPlaceIDs(center, radius, place.UnlimitedCount)
}

func (bf *basicLocalIterator) serveRequest() {
	stopServe := false

	for bf.requests != nil && bf.stop != nil {
		select {
		case c := <-bf.requests:
			for _, placeInfo := range bf.placesInfo {
				c <- &iteratortype.ChargeStationInfo{
					ID: strconv.FormatInt((int64)(placeInfo.ID), 10),
					Location: nav.Location{
						Lat: placeInfo.Location.Lat,
						Lon: placeInfo.Location.Lon,
					},
				}
			}
			close(c)

			if stopServe == true && len(bf.requests) == 0 {
				close(bf.requests)
				bf.requests = nil
			}

		case stopServe = <-bf.stop:
			bf.stop = nil
		}
	}

}

// IterateNearbyStations returns channel contains near by stations
func (bf *basicLocalIterator) IterateNearbyStations() <-chan *iteratortype.ChargeStationInfo {
	c := make(chan *iteratortype.ChargeStationInfo, defaultChargeStaionChannelSize)
	if bf.requests != nil {
		bf.requests <- c
	} else {
		glog.Warning("Call iterator on Stopped Finder, please check your logic.\n")
		close(c)
	}

	return c
}

// Stop stops functionality of finder
func (bf *basicLocalIterator) Stop() {
	if bf.stop != nil {
		bf.stop <- true
	}
}
