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
	rpcMode string
	wayIDs  idsflag.IDs
}

func init() {
	flag.StringVar(&flags.rpcMode, "mode", rpcModeGetWays, "RPC request mode, possible options: "+fmt.Sprintf("%s,%s,%s", rpcModeGetWays, rpcModeGetAll, rpcModeStreamingDelta))
	flag.Var(&flags.wayIDs, "ways", "wayIDs for querying traffic. Use comma-seperated list if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
}
