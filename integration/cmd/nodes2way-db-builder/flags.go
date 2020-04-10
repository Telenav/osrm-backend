package main

import (
	"flag"
)

var flags struct {
	snappyCompressed bool
	in               string // E.g., wayid2nodeids.csv or wayid2nodeids.csv.snappy
	out              string // DB path
}

func init() {
	flag.StringVar(&flags.in, "i", "wayid2nodeids.csv.snappy", "Input wayid2nodeids csv file path, e.g., wayid2nodeids.csv or wayid2nodeids.csv.snappy.")
	flag.StringVar(&flags.out, "o", "nodes2way.db", "Output DB file path.")
	flag.BoolVar(&flags.snappyCompressed, "snappy-compressed", true, "Whether input csv snappy compressed or not.")
}
