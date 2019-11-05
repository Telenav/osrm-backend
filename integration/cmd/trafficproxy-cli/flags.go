package main

import (
	"flag"
	"fmt"
)

const (
	rpcModeGetWays        = "getways"
	rpcModeGetAll         = "getall"
	rpcModeStreamingDelta = "delta"
)

var flags struct {
	rpcMode      string
	wayIDsFlag   wayIDsFlag
	blockingOnly bool
	dumpFile     string
	stdout       bool
}

func init() {
	flag.StringVar(&flags.rpcMode, "mode", rpcModeGetWays, "RPC request mode, possible options: "+fmt.Sprintf("%s,%s,%s", rpcModeGetWays, rpcModeGetAll, rpcModeStreamingDelta))
	flag.Var(&flags.wayIDsFlag, "ways", "wayIDs for querying traffic. Use comma-seperated list if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
	flag.BoolVar(&flags.blockingOnly, "blocking-only", false, "Only use blocking only(blocking flow or blocking incident) live traffic.")
	flag.StringVar(&flags.dumpFile, "dumpfile", "", "Dump file name of flows,incidents. Flows,incident will be dumped to files(xxx_flows.csv,xxx_incidents.csv) if this option is not empty.")
	flag.BoolVar(&flags.stdout, "stdout", true, "Dump flows,incidents to stdout.")
}
