package main

import (
	"flag"
	"fmt"
)

var flags struct {
	pbf       string
	pbfSource string

	historicalSpeedWaysMappingFile       string // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
	outputHistoricalSpeedWaysMappingFile string // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S,TIME_ZONE,DAYLIGHT_SAVING
	withCSVHeader                        bool   // whether dump with CSV header or not
}

const (
	pbfSourceUniDB = "unidb"
)

func init() {
	flag.StringVar(&flags.pbf, "pbf", "", "Input pbf file.")
	flag.StringVar(&flags.pbfSource, "pbf_source", pbfSourceUniDB, fmt.Sprintf("pbf data source, only '%s' is supported since OSM almost no timezone data.", pbfSourceUniDB))
	flag.StringVar(&flags.historicalSpeedWaysMappingFile, "hs-waysmapping-in", "", "Historical speed wayIDs to daily patterns mapping csv file. Pass in multiple files separated by ','.")
	flag.StringVar(&flags.outputHistoricalSpeedWaysMappingFile, "hs-waysmapping-out", "hs-waysmapping.csv", "Historical speed wayIDs to daily patterns, timezone, daylight saving mapping csv file.")
	flag.BoolVar(&flags.withCSVHeader, "csv-header", true, "Whether dump csv header or not.")
}
