package main

import (
	"flag"
	"runtime"
)

var flags struct {
	inputFile           string
	outputFolder        string
	numberOfWorkers     int
	maxRange            float64
	osrmBackendEndpoint string
}

// defaultRange indicates default max drive range in meters.
// This value always bigger than max driven ranges for existing vehicles.
// During query time could just take subset of pre-processed result based on real range.
const defaultMaxRange = 800000

func init() {
	flag.StringVar(&flags.inputFile, "i", "", "path for input file in json format")
	flag.StringVar(&flags.outputFolder, "o", "", "path for output folder")
	flag.IntVar(&flags.numberOfWorkers, "num_of_workers", runtime.NumCPU(), "number of workers to build connectivity map")
	flag.Float64Var(&flags.maxRange, "range", defaultMaxRange, "maximum drive range in meters used for preprocessing.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "OSRM-backend endpoint")
}
