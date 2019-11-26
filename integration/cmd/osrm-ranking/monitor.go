package main

import (
	"fmt"
	"time"
)

type monitorContents struct {
	UpTime                       jsonDuration                  `json:"uptime"`
	TrafficCacheMonitorContents  *trafficCacheMonitorContents  `json:"traffic_cache"`
	WayID2NodeIDsMonitorContents *wayID2NodeIDsMonitorContents `json:"wayids2nodeids"`
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

func newMonitorContents() *monitorContents {
	return &monitorContents{
		0, &trafficCacheMonitorContents{}, &wayID2NodeIDsMonitorContents{},
	}
}

type jsonDuration time.Duration

func (j jsonDuration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Duration(j).String())), nil
}
