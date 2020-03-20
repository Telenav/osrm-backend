package main

import (
	"strings"

	"github.com/Telenav/osrm-backend/integration/traffic/historicalspeed"
)

func newWaysMappingUpdater(inFiles []string, outFile string, withCSVHeader bool, inWayTimezoneInfo <-chan *wayTimezoneInfo) error {

	waysMapping := historicalspeed.WaysMapping{}
	if err := waysMapping.Load(strings.Split(flags.historicalSpeedWaysMappingFile, ",")); err != nil {
		return err
	}

	for {
		_, ok := <-inWayTimezoneInfo
		if !ok {
			break
		}

		//TODO: update way's timezone and daylight saving in waysmapping
	}

	if err := waysMapping.Dump(outFile, withCSVHeader); err != nil {
		return err
	}

	return nil
}
