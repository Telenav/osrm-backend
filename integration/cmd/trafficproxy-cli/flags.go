package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/wayidsflag"
)

const (
	rpcModeGetWays        = "getways"
	rpcModeGetAll         = "getall"
	rpcModeStreamingDelta = "delta"
)

var flags struct {
	rpcMode                      string
	wayIDs                       wayidsflag.WayIDs
	blockingOnly                 bool
	dumpFile                     string
	stdout                       bool
	humanFriendlyCSV             bool
	streamingDeltaDumpInterval   time.Duration
	streamingDeltaSplitDumpFiles bool
}

func init() {
	flag.StringVar(&flags.rpcMode, "mode", rpcModeGetWays, "RPC request mode, possible options: "+fmt.Sprintf("%s,%s,%s", rpcModeGetWays, rpcModeGetAll, rpcModeStreamingDelta))
	flag.Var(&flags.wayIDs, "ways", "wayIDs for querying traffic. Use comma-seperated list if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
	flag.BoolVar(&flags.blockingOnly, "blocking-only", false, "Only use blocking only(blocking flow or blocking incident) live traffic.")
	flag.StringVar(&flags.dumpFile, "dumpfile", "", "Dump file name of flows,incidents. Flows,incident will be dumped to files(xxx_flows.csv,xxx_incidents.csv) if this option is not empty.")
	flag.BoolVar(&flags.stdout, "stdout", true, "Dump flows,incidents to stdout.")
	flag.BoolVar(&flags.humanFriendlyCSV, "humanfriendly", false, "Human friendly contents in csv, i.e. prefer string instead of integer/boolean as much as possible in csv files. E.g. TrafficLevel, IncidentType, IncidentSeverity, IsBlocking.")
	flag.DurationVar(&flags.streamingDeltaDumpInterval, "delta-dump-interval", 60*time.Second, "Dump streaming delta traffic flows,incidents interval, e.g. split dump files, statistics, etc.")
	flag.BoolVar(&flags.streamingDeltaSplitDumpFiles, "delta-dump-split", true, "Whether split dump files per delta-dump-interval.")
}
