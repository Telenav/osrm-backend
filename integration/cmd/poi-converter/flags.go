package main

import "flag"

var flags struct {
	inputPath  string
	outputPath string
}

func init() {
	flag.StringVar(&flags.inputPath, "i", "", "path for input file in csv format")
	flag.StringVar(&flags.outputPath, "o", "output.json", "path for output file in json format")
}
