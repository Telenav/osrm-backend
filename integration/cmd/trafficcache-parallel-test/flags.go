package main

import (
	"flag"
	"time"
)

var flags struct {
	monitorInterval          time.Duration
	wayID2NodeIDsMappingFile string
	indexedByWayID           bool
	indexedByEdge            bool
}

func init() {
	flag.DurationVar(&flags.monitorInterval, "monitor-interval", 10*time.Second, "Log for traffic cache status will print out per monitor-interval.")
	flag.StringVar(&flags.wayID2NodeIDsMappingFile, "m", "wayid2nodeids.csv.snappy", "OSRM way id to node ids mapping table, snappy compressed.")
	flag.BoolVar(&flags.indexedByWayID, "indexed-by-way", true, "Run cache indexed by wayID.")
	flag.BoolVar(&flags.indexedByEdge, "indexed-by-edge", true, "Run cache indexed by Edge.")
}
