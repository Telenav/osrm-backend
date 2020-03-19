package main

import (
	"flag"
	"fmt"
)

var flags struct {
	input     string
	output    string
	pbfSource string
}

const (
	pbfSourceUniDB = "unidb"
	pbfSourceOSM   = "osm"
)

func init() {
	flag.StringVar(&flags.input, "i", "", "Input pbf file.")
	flag.StringVar(&flags.output, "o", "", "Output csv file")
	flag.StringVar(&flags.pbfSource, "pbf_source", pbfSourceOSM, fmt.Sprintf("pbf data source, can be '%s' or '%s'.", pbfSourceOSM, pbfSourceUniDB))
}
