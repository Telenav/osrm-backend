package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/rankingservice"
	"github.com/Telenav/osrm-backend/integration/trafficcache/trafficcacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/wayid2nodeids"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// monitor
	upClock := time.Now()
	monitorContents := newMonitorContents()
	monitorContents.TrafficCacheMonitorContents.Name = "traffic cache(indexed by edge)"

	// wayid2nodeids mapping
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

	//start http listening
	mux := http.NewServeMux()
	mux.HandleFunc("/monitor/", func(w http.ResponseWriter, req *http.Request) {
		monitorContents.UpTime = jsonDuration(time.Now().Sub(upClock))

		// update wayid2nodeids contents
		monitorContents.WayID2NodeIDsMonitorContents.IsReady = wayID2NodeIDsMapping.IsReady()
		monitorContents.WayID2NodeIDsMonitorContents.Ways = wayID2NodeIDsMapping.WayIDsCount()

		// update traffic cache contents
		monitorContents.TrafficCacheMonitorContents.Flows = trafficCache.Flows.Count()
		monitorContents.TrafficCacheMonitorContents.FlowsAffectedWays = trafficCache.Flows.AffectedWaysCount()
		monitorContents.TrafficCacheMonitorContents.Incidents = trafficCache.Incidents.Count()
		monitorContents.TrafficCacheMonitorContents.IncidentsAffectedWays = trafficCache.Incidents.AffectedWaysCount()
		monitorContents.TrafficCacheMonitorContents.IncidentsAffectedEdges = trafficCache.Incidents.AffectedEdgesCount()
		glog.Infof("monitor %s, [flows] %d affectedways %d, [incidents] blocking-only %d, affectedways %d affectededges %d",
			monitorContents.TrafficCacheMonitorContents.Name, monitorContents.TrafficCacheMonitorContents.Flows, monitorContents.TrafficCacheMonitorContents.FlowsAffectedWays,
			monitorContents.TrafficCacheMonitorContents.Incidents, monitorContents.TrafficCacheMonitorContents.IncidentsAffectedWays, monitorContents.TrafficCacheMonitorContents.IncidentsAffectedEdges)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(monitorContents)
	})

	//start ranking service
	rankingService := rankingservice.New(flags.osrmBackendEndpoint, trafficCache)
	mux.Handle("/route/v1/driving/", rankingService)

	// listen
	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
