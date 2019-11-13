package main

import (
	"flag"

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
	feeder.Run() //TODO: should async run

	//TODO: start ranking service
}
