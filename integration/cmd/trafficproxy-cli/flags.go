package main

import (
	"flag"
)

var flags struct {
	streamingDelta bool
	wayIDsFlag     wayIDsFlag
	blockingOnly   bool
	dumpFile       string
}

func init() {
	flag.BoolVar(&flags.streamingDelta, "delta", false, "Whether receive traffic data by streaming delta mode. Defaultly false means get traffic data once.")
	flag.Var(&flags.wayIDsFlag, "ways", "wayIDs for querying traffic. Use ',' split if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
	flag.BoolVar(&flags.blockingOnly, "blocking-only", false, "Only use blocking only(blocking flow or blocking incident) live traffic.")
	flag.StringVar(&flags.dumpFile, "dumpfile", "", "Dump file name of flows,incidents. Flows,incident will be dumped to files(xxx_flows.csv,xxx_incidents.csv) instead of stdout.")
}
