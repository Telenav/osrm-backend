package main

import (
	"flag"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/util/mapsource"
)

var flags struct {
	pbf       string
	mapSource string

	inWaysMappingFile  string // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
	outWaysMappingFile string // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S,TIME_ZONE,DAYLIGHT_SAVING
	withCSVHeader      bool   // whether dump with CSV header or not
}

func init() {
	flag.StringVar(&flags.pbf, "pbf", "", "Input pbf file.")
	flag.StringVar(&flags.mapSource, "mapsource", mapsource.UniDB, fmt.Sprintf("pbf map data source, only '%s' is supported since '%s' almost no timezone data.", mapsource.UniDB, mapsource.OSM))
	flag.StringVar(&flags.inWaysMappingFile, "hs-waysmapping-in", "", "Historical speed wayIDs to daily patterns mapping csv file. Pass in multiple files separated by ','.")
	flag.StringVar(&flags.outWaysMappingFile, "hs-waysmapping-out", "hs-waysmapping.csv", "Historical speed wayIDs to daily patterns, timezone, daylight saving mapping csv file.")
	flag.BoolVar(&flags.withCSVHeader, "csv-header", true, "Whether dump csv header or not.")
}
