package main

import (
	"flag"
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/trafficcache"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficcache/trafficcacheindexedbyedge"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/wayid2nodeids"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	var cacheByWay *trafficcache.Cache
	var cacheByEdge *trafficcacheindexedbyedge.Cache
	if flags.indexedByEdge {
		wayID2NodeIDsMapping := wayid2nodeids.NewMappingFrom(flags.wayID2NodeIDsMappingFile)
		if err := wayID2NodeIDsMapping.Load(); err != nil {
			glog.Error(err)
			return
		}
		cacheByEdge = trafficcacheindexedbyedge.New(wayID2NodeIDsMapping)
	}
	if flags.indexedByWayID {
		cacheByWay = trafficcache.New()
	}

	// traffic cache feeder
	feeder := trafficproxyclient.NewFeeder()
	if cacheByWay != nil {
		feeder.RegisterEaters(cacheByWay)
	}
	if cacheByEdge != nil {
		feeder.RegisterEaters(cacheByEdge)
	}

	go func() {
		for {
			err := feeder.Run()
			if err != nil {
				glog.Warning(err)
			}

			if cacheByWay != nil {
				cacheByWay.Clear()
			}
			if cacheByEdge != nil {
				cacheByEdge.Clear()
			}
			time.Sleep(5 * time.Second) // try again later
		}
	}()

	// monitor
	startTime := time.Now()
	for {
		currentTime := time.Now()
		if currentTime.Sub(startTime) < flags.monitorInterval {
			time.Sleep(time.Second)
			continue
		}
		startTime = currentTime

		var cacheByWayFlowsCount, cacheByWayIncidentsCount, cacheByWayIncidentAffectedWaysCount int64
		var cacheByEdgeFlowsAffectedWaysCount, cacheByEdgeIncidentsCount, cacheByEdgeIncidentAffectedWaysCount int64
		if cacheByWay != nil {
			cacheByWayFlowsCount = cacheByWay.Flows.Count()
			cacheByWayIncidentsCount = int64(cacheByWay.Incidents.Count())
			cacheByWayIncidentAffectedWaysCount = int64(cacheByWay.Incidents.AffectedWaysCount())
			glog.Infof("traffic in cache(indexed by wayID), [flows] %d, [incidents] blocking-only %d, affectedways %d",
				cacheByWayFlowsCount, cacheByWayIncidentsCount, cacheByWayIncidentAffectedWaysCount)
		}
		if cacheByEdge != nil {
			cacheByEdgeFlowsAffectedWaysCount = cacheByEdge.Flows.AffectedWaysCount()
			cacheByEdgeIncidentsCount = int64(cacheByEdge.Incidents.Count())
			cacheByEdgeIncidentAffectedWaysCount = int64(cacheByEdge.Incidents.AffectedWaysCount())
			glog.Infof("traffic in cache(indexed by Edge), [flows] %d affectedways %d, [incidents] blocking-only %d, affectedways %d affectededges %d",
				cacheByEdge.Flows.Count(), cacheByEdgeFlowsAffectedWaysCount,
				cacheByEdgeIncidentsCount, cacheByEdgeIncidentAffectedWaysCount, cacheByEdge.Incidents.AffectedEdgesCount())
		}
		if cacheByWay != nil && cacheByEdge != nil {
			warnMismatch("flows", cacheByWayFlowsCount, cacheByEdgeFlowsAffectedWaysCount)
			warnMismatch("incidents", cacheByWayIncidentsCount, cacheByEdgeIncidentsCount)
			warnMismatch("incidents affected ways", cacheByWayIncidentAffectedWaysCount, cacheByEdgeIncidentAffectedWaysCount)
		}
	}
}

func warnMismatch(name string, v1, v2 int64) {
	if v1 != v2 {
		glog.Warningf("%s mismatch: %d != %d, delta %d", name, v1, v2, v1-v2)
	}
}
