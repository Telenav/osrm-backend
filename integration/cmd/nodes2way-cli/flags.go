package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/pkg/wayidsflag"
)

var flags struct {
	nodeIDs wayidsflag.WayIDs //TODO: will refactor it to intsflag.Int64s later
	db      string
	dbStat  bool
}

func init() {
	flag.Var(&flags.nodeIDs, "nodes", "Continuously comma-seperated nodeIDs on a route. E.g. '167772220006101,167772220007101,167772220008101'.")
	flag.StringVar(&flags.db, "db", "nodes2way.db", "DB file path.")
	flag.BoolVar(&flags.dbStat, "dbstat", false, "Print DB statistics.")
}
