package main

import (
	"flag"
)

var flags struct {
	streamingDelta bool
	wayIDsFlag     wayIDsFlag
}

func init() {
	flag.BoolVar(&flags.streamingDelta, "delta", false, "Whether receive traffic data by streaming delta mode. Defaultly false means get traffic data once.")
	flag.Var(&flags.wayIDsFlag, "ways", "wayIDs for querying traffic. Use ',' split if more than one wayID. Positive value means forward, negative value means backward. E.g. '829733412,-104489539'.")
}
