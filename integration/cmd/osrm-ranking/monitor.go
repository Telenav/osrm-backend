package main

import (
	"fmt"
	"time"
)

type monitorContents struct {
	UpTime                         jsonDuration                    `json:"uptime"`
	HistoricalSpeedMonitorContents *historicalSpeedMonitorContents `json:"historical_speed"`
	TrafficCacheMonitorContents    *trafficCacheMonitorContents    `json:"traffic_cache"`
	WayID2NodeIDsMonitorContents   *wayID2NodeIDsMonitorContents   `json:"wayid2nodeids"`
	CmdlineArgs                    []string                        `json:"cmdline"`
}

type trafficCacheMonitorContents struct {
	Name                   string `json:"name"`
	Flows                  int64  `json:"flows"`
	FlowsAffectedWays      int64  `json:"flows_affected_ways"`
	Incidents              int    `json:"incidents"`
	IncidentsAffectedWays  int    `json:"incidents_affected_ways"`
	IncidentsAffectedEdges int    `json:"incidents_affected_edges"`
}

type wayID2NodeIDsMonitorContents struct {
	IsReady bool `json:"is_ready"`
	Ways    int  `json:"ways"`
}

type historicalSpeedMonitorContents struct {
	DailyPatterns       int `json:"daily_patterns"`
	Way2PatternsMapping int `json:"way2patterns"`
}

func newMonitorContents() *monitorContents {
	return &monitorContents{
		0, &historicalSpeedMonitorContents{}, &trafficCacheMonitorContents{}, &wayID2NodeIDsMonitorContents{}, nil,
	}
}

type jsonDuration time.Duration

func (j jsonDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(j).String())), nil
}
