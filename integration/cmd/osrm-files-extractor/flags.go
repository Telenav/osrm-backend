package main

import "flag"

var flags struct {
	filePath     string
	printSummary int
}

func init() {
	flag.StringVar(&flags.filePath, "f", "", "Single OSRM file to load, e.g. 'nevada-latest.osrm' or 'nevada-latest.osrm.nbg_nodes'.")
	flag.IntVar(&flags.printSummary, "summary", -1, "Print summary and head lines of loaded contents. <0: not print; ==0: only print summary; >0: print summary and head lines.")
}
