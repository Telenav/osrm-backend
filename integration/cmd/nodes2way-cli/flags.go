package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/pkg/wayidsflag"
)

var flags struct {
	nodeIDs wayidsflag.WayIDs //TODO: will refactor it to intsflag.Int64s later
	db      string
}

func init() {
	flag.StringVar(&flags.db, "db", "nodes2way.db", "DB file path.")
	flag.Var(&flags.nodeIDs, "nodes", "Continuously comma-seperated nodeIDs on a route. E.g. '829733412,104489539'.")
}
