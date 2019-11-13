package main

import (
	"flag"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/trafficnodepaircache"
	"github.com/Telenav/osrm-backend/integration/trafficwayidcache"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// prepare traffic cache
	cacheIndexedByWayID := trafficwayidcache.New()
	cacheIndexedByNodePair := trafficnodepaircache.New()
	feeder := trafficproxyclient.NewFeeder()
	feeder.RegisterEaters(cacheIndexedByWayID, cacheIndexedByNodePair)
	for { //TODO: should async run
		err := feeder.Run()
		if err != nil {
			glog.Warning(err)
		}
		cacheIndexedByNodePair.Clear()
		cacheIndexedByWayID.Clear()
		time.Sleep(5 * time.Second) // try again later
	}

	//TODO: start ranking service
}
