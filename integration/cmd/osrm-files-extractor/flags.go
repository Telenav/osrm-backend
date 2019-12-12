package main

import "flag"

var flags struct {
	filePath     string
	printSummary bool
}

func init() {
	flag.StringVar(&flags.filePath, "f", "", "Single OSRM file to load, e.g. 'nevada-latest.osrm' or 'nevada-latest.osrm.nbg_nodes'.")
	flag.BoolVar(&flags.printSummary, "summary", false, "Print summary of loaded contents.")
}
