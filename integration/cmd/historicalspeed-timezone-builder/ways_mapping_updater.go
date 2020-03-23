package main

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/historicalspeed"
	"github.com/golang/glog"
)

func newWaysMappingUpdater(inFiles []string, outFile string, withCSVHeader bool, inWayTimezoneInfo <-chan *wayTimezoneInfo) error {

	waysMapping := historicalspeed.WaysMapping{}
	if err := waysMapping.Load(inFiles); err != nil {
		return err
	}

	startTime := time.Now()
	var inCount, updateSucceedcount int
	for {
		tzInfo, ok := <-inWayTimezoneInfo
		if !ok {
			break
		}
		inCount++

		if err := waysMapping.UpdateTimezoneDaylightSaving(tzInfo.wayID, tzInfo.timezone, tzInfo.daylightSaving); err != nil {
			if glog.V(3) { // avoid affect performance by verbose log
				glog.Infof("Update timezone info %+v for historical speed failed, err: %v", tzInfo, err)
			}
			continue
		}
		updateSucceedcount++
	}
	glog.Infof("Updated timezone/daylight saving for historical speed ways mapping, total income %d, update succeed %d, takes %f seconds", inCount, updateSucceedcount, time.Now().Sub(startTime).Seconds())

	if err := waysMapping.Dump(outFile, withCSVHeader); err != nil {
		return err
	}

	return nil
}
