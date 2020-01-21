package main

import "flag"

var flags struct {
	filePath     string
	singleFile   bool
	printSummary int
}

func init() {
	flag.StringVar(&flags.filePath, "f", "", "OSRM files(or a single file) to load, e.g. 'nevada-latest.osrm' or 'nevada-latest.osrm.nbg_nodes'. If input is 'xxx.osrm', depends on '-single_file' to load it only or load all .osrm.xxx.")
	flag.BoolVar(&flags.singleFile, "single_file", false, "Only valid if the file path is 'xxx.osrm' from '-f'. false to load all xxx.osrm.xxx automatically, true to load the single xxx.osrm file only.")
	flag.IntVar(&flags.printSummary, "summary", -1, "Print summary and head lines of loaded contents. <0: not print; ==0: only print summary; >0: print summary and head lines.")
}
