package localfinder

import (
	"strconv"

	"github.com/Telenav/osrm-backend/integration/api/nav"
	"github.com/Telenav/osrm-backend/integration/service/oasis/internal/entity"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/service/oasis/place/stationfinder/stationfindertype"
	"github.com/golang/glog"
)

const defaultIteratorCount = 5
const defaultChargeStaionChannelSize = 500

type basicLocalFinder struct {
	localFinder place.Finder
	placesInfo  []*entity.PlaceInfo
	requests    chan chan *stationfindertype.ChargeStationInfo
	stop        chan bool
}

func newBasicLocalFinder(localFinder place.Finder) *basicLocalFinder {
	bf := &basicLocalFinder{
		localFinder: localFinder,
		requests:    make(chan chan *stationfindertype.ChargeStationInfo, defaultIteratorCount),
		stop:        make(chan bool),
	}
	go bf.serveRequest()
	return bf
}

func (bf *basicLocalFinder) getNearbyChargeStations(center nav.Location, radius float64) {
	bf.placesInfo = bf.localFinder.FindNearByPlaceIDs(center, radius, place.UnlimitedCount)
}

func (bf *basicLocalFinder) serveRequest() {
	stopServe := false

	for bf.requests != nil && bf.stop != nil {
		select {
		case c := <-bf.requests:
			for _, placeInfo := range bf.placesInfo {
				c <- &stationfindertype.ChargeStationInfo{
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
func (bf *basicLocalFinder) IterateNearbyStations() <-chan *stationfindertype.ChargeStationInfo {
	c := make(chan *stationfindertype.ChargeStationInfo, defaultChargeStaionChannelSize)
	if bf.requests != nil {
		bf.requests <- c
	} else {
		glog.Warning("Call iterator on Stopped Finder, please check your logic.\n")
		close(c)
	}

	return c
}

// Stop stops functionality of finder
func (bf *basicLocalFinder) Stop() {
	if bf.stop != nil {
		bf.stop <- true
	}
}
