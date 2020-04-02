package trafficdumper

import (
	"flag"
	"time"
)

var flags struct {
	blockingOnly                 bool
	dumpFile                     string
	stdout                       bool
	humanFriendlyCSV             bool
	streamingDeltaDumpInterval   time.Duration
	streamingDeltaSplitDumpFiles bool
}

func init() {
	flag.BoolVar(&flags.blockingOnly, "blocking-only", false, "Only use blocking only(blocking flow or blocking incident) live traffic.")
	flag.StringVar(&flags.dumpFile, "dumpfile", "", "Dump file name of flows,incidents. Flows,incident will be dumped to files(xxx_flows.csv,xxx_incidents.csv) if this option is not empty.")
	flag.BoolVar(&flags.stdout, "stdout", true, "Dump flows,incidents to stdout.")
	flag.BoolVar(&flags.humanFriendlyCSV, "humanfriendly", false, "Human friendly contents in csv, i.e. prefer string instead of integer/boolean as much as possible in csv files. E.g. TrafficLevel, IncidentType, IncidentSeverity, IsBlocking.")
	flag.DurationVar(&flags.streamingDeltaDumpInterval, "delta-dump-interval", 60*time.Second, "Dump streaming delta traffic flows,incidents interval, e.g. split dump files, statistics, etc.")
	flag.BoolVar(&flags.streamingDeltaSplitDumpFiles, "delta-dump-split", true, "Whether split dump files per delta-dump-interval.")
}
