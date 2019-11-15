package main

import (
	"flag"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/trafficcacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/trafficcacheindexedbywayid"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// prepare traffic cache
	cacheByWay := trafficcacheindexedbywayid.New()
	cacheByEdge := trafficcacheindexedbyedge.New()
	feeder := trafficproxyclient.NewFeeder()
	feeder.RegisterEaters(cacheByWay, cacheByEdge)
	for { //TODO: should async run
		err := feeder.Run()
		if err != nil {
			glog.Warning(err)
		}
		cacheByWay.Clear()
		cacheByEdge.Clear()
		time.Sleep(5 * time.Second) // try again later
	}

	//TODO: start ranking service
}
