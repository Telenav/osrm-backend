package main

import (
	"flag"
	"time"
)

var flags struct {
	listenPort               int
	monitorInterval          time.Duration
	wayID2NodeIDsMappingFile string
	osrmBackendEndpoint      string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8080, "Listen port.")
	flag.DurationVar(&flags.monitorInterval, "monitor-interval", 10*time.Second, "Log for traffic cache status will print out per monitor-interval.")
	flag.StringVar(&flags.wayID2NodeIDsMappingFile, "m", "wayid2nodeids.csv.snappy", "OSRM way id to node ids mapping table, snappy compressed.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "Backend OSRM-backend endpoint")
}
