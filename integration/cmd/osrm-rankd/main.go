package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Telenav/osrm-backend/integration/service/ranking"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel"
	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/modelfactory"
	"github.com/Telenav/osrm-backend/integration/traffic/historicalspeed"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes/nodes2wayblotdb"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if flags.version {
		printVersion()
		return
	}

	// monitor
	upClock := time.Now()
	monitorContents := newMonitorContents()
	monitorContents.Versions.AppVersion = appVersion
	monitorContents.Versions.GitCommit = gitCommit
	monitorContents.Versions.BuildTime = buildTime
	monitorContents.CmdlineArgs = os.Args

	// prepare nodes2way
	nodes2wayDB, err := nodes2wayblotdb.Open(flags.nodes2WayDB, true)
	if err != nil {
		glog.Errorf("Load nodes2way DB failed, err: %v", err)
		return
	}
	monitorContents.Nodes2WayDB = nodes2wayDB.Statistics()

	// prepare historical speeds if available
	var historicalSpeedCache *historicalspeed.Speeds
	if flags.historicalSpeed {
		historicalSpeedCache = historicalspeed.New(strings.Split(flags.historicalSpeedDailyPatternFile, ","), strings.Split(flags.historicalSpeedWaysMappingFile, ","))
		if err := historicalSpeedCache.Load(); err != nil {
			glog.Errorf("Load historical speed failed, err: %v", err)
			return
		}
		glog.Infof("Historical speeds loaded: daily patterns count %d, ways(directed) count %d.", historicalSpeedCache.DailyPatternsCount(), historicalSpeedCache.WaysCount())
		monitorContents.HistoricalSpeedMonitorContents.DailyPatterns = historicalSpeedCache.DailyPatternsCount()
		monitorContents.HistoricalSpeedMonitorContents.Way2PatternsMapping = historicalSpeedCache.WaysCount()
	}

	// prepare traffic cache
	var liveTrafficCache *trafficcache.Cache
	if flags.liveTraffic {
		liveTrafficCache = trafficcache.New()
		feeder := trafficproxyclient.NewFeeder()
		feeder.RegisterEaters(liveTrafficCache)
		go func() {
			for {
				err := feeder.Run()
				if err != nil {
					glog.Warning(err)
				}
				liveTrafficCache.Clear()
				time.Sleep(5 * time.Second) // try again later
			}
		}()
	}

	//start http listening
	mux := http.NewServeMux()

	//monitor service
	mux.HandleFunc("/monitor/", func(w http.ResponseWriter, req *http.Request) {
		monitorContents.UpTime = jsonDuration(time.Now().Sub(upClock))

		// update traffic cache contents
		if liveTrafficCache != nil {
			monitorContents.TrafficCacheMonitorContents.Flows = liveTrafficCache.Flows.Count()
			monitorContents.TrafficCacheMonitorContents.Incidents = liveTrafficCache.Incidents.Count()
			monitorContents.TrafficCacheMonitorContents.IncidentsAffectedWays = liveTrafficCache.Incidents.AffectedWaysCount()
			glog.Infof("monitor live traffic, [flows] %d, [incidents] blocking-only %d, affectedways %d.",
				monitorContents.TrafficCacheMonitorContents.Flows,
				monitorContents.TrafficCacheMonitorContents.Incidents, monitorContents.TrafficCacheMonitorContents.IncidentsAffectedWays)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(monitorContents)
	})

	//start ranking service
	var trafficApplier trafficapplyingmodel.Applier
	if liveTrafficCache != nil || historicalSpeedCache != nil {
		var err error
		trafficApplier, err = modelfactory.NewApplier(flags.trafficApplyingModel, liveTrafficCache, historicalSpeedCache)
		if err != nil {
			glog.Errorf("New traffic applying model failed, err %v", err)
			return
		}
	}
	rankingService := ranking.New(flags.osrmBackendEndpoint, nodes2wayDB, trafficApplier)
	mux.Handle("/route/v1/driving/", rankingService)

	// listen
	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
