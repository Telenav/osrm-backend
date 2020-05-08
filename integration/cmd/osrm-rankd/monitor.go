package main

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/appversion"
)

type monitorContents struct {
	UpTime                         jsonDuration `json:"uptime"`
	appversion.VersionInfo         `json:"version"`
	HistoricalSpeedMonitorContents *historicalSpeedMonitorContents `json:"historical speed"`
	TrafficCacheMonitorContents    *trafficCacheMonitorContents    `json:"live traffic"`
	Nodes2WayDB                    string                          `json:"nodes2way"`
	CmdlineArgs                    []string                        `json:"cmdline"`
}

type trafficCacheMonitorContents struct {
	Flows                 int64 `json:"flows"`
	Incidents             int   `json:"block_incidents"`
	IncidentsAffectedWays int   `json:"incidents_affected_ways"`
}

type historicalSpeedMonitorContents struct {
	DailyPatterns       int `json:"daily_patterns"`
	Way2PatternsMapping int `json:"way2patterns"`
}

func newMonitorContents() *monitorContents {
	return &monitorContents{
		0, appversion.VersionInfo{}, &historicalSpeedMonitorContents{}, &trafficCacheMonitorContents{}, "", nil,
	}
}

type jsonDuration time.Duration

func (j jsonDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(j).String())), nil
}
