package main

import (
	"flag"
)

var flags struct {
	listenPort               int
	wayID2NodeIDsMappingFile string
	osrmBackendEndpoint      string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8080, "Listen port.")
	flag.StringVar(&flags.wayID2NodeIDsMappingFile, "m", "wayid2nodeids.csv.snappy", "OSRM way id to node ids mapping table, snappy compressed.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "Backend OSRM-backend endpoint")
}
