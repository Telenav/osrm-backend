package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/backend"
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/ranking"
	"github.com/Telenav/osrm-backend/integration/trafficcache/trafficcacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/wayid2nodeids"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	wayID2NodeIDsMapping := wayid2nodeids.NewMappingFrom(flags.wayID2NodeIDsMappingFile)
	if err := wayID2NodeIDsMapping.Load(); err != nil {
		glog.Error(err)
		return
	}

	// prepare traffic cache
	trafficCache := trafficcacheindexedbyedge.New(wayID2NodeIDsMapping)
	feeder := trafficproxyclient.NewFeeder()
	feeder.RegisterEaters(trafficCache)
	go func() {
		for {
			err := feeder.Run()
			if err != nil {
				glog.Warning(err)
			}
			trafficCache.Clear()
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

			glog.Infof("traffic in cache(indexed by Edge), [flows] %d affectedways %d, [incidents] blocking-only %d, affectedways %d affectededges %d",
				trafficCache.Flows.Count(), trafficCache.Flows.AffectedWaysCount(),
				trafficCache.Incidents.Count(), trafficCache.Incidents.AffectedWaysCount(), trafficCache.Incidents.AffectedEdgesCount())
		}
	}()

	//start http listening
	mux := http.NewServeMux()

	mux.HandleFunc("/monitor/", func(w http.ResponseWriter, req *http.Request) {
		//TODO:

		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "Not implemented")
	})

	//start ranking service
	rankingService := ranking.New(flags.osrmBackendEndpoint, backend.Timeout(), trafficCache)
	mux.Handle("/route/v1/driving/", rankingService)

	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
