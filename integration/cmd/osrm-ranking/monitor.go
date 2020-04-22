package main

import (
	"fmt"
	"time"
)

type monitorContents struct {
	UpTime                         jsonDuration                    `json:"uptime"`
	HistoricalSpeedMonitorContents *historicalSpeedMonitorContents `json:"historical speed"`
	TrafficCacheMonitorContents    *trafficCacheMonitorContents    `json:"live traffic"`
	Nodes2WayDB                    string                          `json:"nodes2way"`
	CmdlineArgs                    []string                        `json:"cmdline"`
}

type trafficCacheMonitorContents struct {
	Flows                  int64 `json:"flows"`
	Incidents              int   `json:"block_incidents"`
	IncidentsAffectedWays  int   `json:"incidents_affected_ways"`
	IncidentsAffectedEdges int   `json:"incidents_affected_edges"`
}

type historicalSpeedMonitorContents struct {
	DailyPatterns       int `json:"daily_patterns"`
	Way2PatternsMapping int `json:"way2patterns"`
}

func newMonitorContents() *monitorContents {
	return &monitorContents{
		0, &historicalSpeedMonitorContents{}, &trafficCacheMonitorContents{}, "", nil,
	}
}

type jsonDuration time.Duration

func (j jsonDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(j).String())), nil
}
