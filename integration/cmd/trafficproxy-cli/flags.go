package main

import (
	"flag"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/util/idsflag"
)

const (
	rpcModeGetWays        = "getways"
	rpcModeGetAll         = "getall"
	rpcModeStreamingDelta = "delta"
)

var flags struct {
	version bool // print version

	rpcMode string
	wayIDs  idsflag.IDs
}

func init() {
	flag.BoolVar(&flags.version, "version", false, "Print version and exit.")

	flag.StringVar(&flags.rpcMode, "mode", rpcModeGetWays, "RPC request mode, possible options: "+fmt.Sprintf("%s,%s,%s", rpcModeGetWays, rpcModeGetAll, rpcModeStreamingDelta))
	flag.Var(&flags.wayIDs, "ways", "wayIDs for querying traffic. Use comma-seperated list if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
}
