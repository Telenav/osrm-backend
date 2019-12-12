package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotosrm"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	contents, err := dotosrm.Load(flags.filePath)
	if err != nil {
		glog.Error(err)
		return
	}
	if flags.printSummary {
		contents.PrintSummary()
	}
}
