package main

import "flag"

var flags struct {
	filePath     string
	singleFile   bool
	printSummary int
	packBits     uint // https://github.com/Telenav/osrm-backend/blob/6283c6074066f98e6d4a9f774f21ea45407c0d52/include/extractor/packed_osm_ids.hpp#L14
}

func init() {
	flag.StringVar(&flags.filePath, "f", "", "OSRM files(or a single file) to load, e.g. 'nevada-latest.osrm' or 'nevada-latest.osrm.nbg_nodes'. If input is 'xxx.osrm', depends on '-single_file' to load it only or load all .osrm.xxx.")
	flag.BoolVar(&flags.singleFile, "single_file", false, "Only valid if the file path is 'xxx.osrm' from '-f'. false to load all xxx.osrm.xxx automatically, true to load the single xxx.osrm file only.")
	flag.IntVar(&flags.printSummary, "summary", -1, "Print summary and head lines of loaded contents. <0: not print; ==0: only print summary; >0: print summary and head lines.")
	flag.UintVar(&flags.packBits, "packed_bits", 63, "Bits for parsing packed_vector PackedOSMIDs, range [1,64]. It's 33 bits in osrm/osrm-backend, which Telenav/osrm-backend uses 63 bits instead.")
}
