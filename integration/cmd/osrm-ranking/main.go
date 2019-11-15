package main

import (
	"flag"
	"net/http"
	"strconv"
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
	go func() {
		for {
			err := feeder.Run()
			if err != nil {
				glog.Warning(err)
			}
			cacheByWay.Clear()
			cacheByEdge.Clear()
			time.Sleep(5 * time.Second) // try again later
		}
	}()

	// monitor
	go func() {

		startTime := time.Now()
		for {
			currentTime := time.Now()
			if currentTime.Sub(startTime) < flags.monitorInterval {
				time.Sleep(time.Second)
				continue
			}
			startTime = currentTime

			if cacheByWay != nil {
				incidents, waysAffectedByIncidents := cacheByWay.IncidentsAndAffectedWaysCount()
				glog.Infof("traffic in cache(indexed by wayID), flows: %d, incidents(blocking-only): %d, ways(affected by incidents): %d",
					cacheByWay.FlowCount(), incidents, waysAffectedByIncidents)
			}
			if cacheByEdge != nil {
				//TODO:
			}
		}
	}()

	//start ranking service
	mux := http.NewServeMux()
	//TODO:

	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
