package main

import (
	"flag"
)

var flags struct {
	listenPort int

	wayID2NodeIDsMappingFile string

	historicalSpeed                        bool // whether enable historical speed or not
	historicalSpeedDailyPatternFile        string
	historicalSpeedWay2PatternsMappingFile string

	osrmBackendEndpoint string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8080, "Listen port.")

	flag.StringVar(&flags.wayID2NodeIDsMappingFile, "m", "wayid2nodeids.csv.snappy", "OSRM way id to node ids mapping table, snappy compressed.")

	flag.BoolVar(&flags.historicalSpeed, "hs", false, "Enable historical speed. The historical speed related files won't be loaded if disabled.")
	flag.StringVar(&flags.historicalSpeedDailyPatternFile, "hs-dailypattern", "", "Historical speed daily patterns csv file.")
	flag.StringVar(&flags.historicalSpeedWay2PatternsMappingFile, "hs-wayid2patterns", "", "Historical speed wayIDs to daily patterns mapping csv file. Pass in multiple files separated by ','.")

	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "Backend OSRM-backend endpoint")
}
