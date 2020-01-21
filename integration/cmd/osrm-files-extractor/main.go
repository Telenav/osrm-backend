package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotosrm"
	"github.com/Telenav/osrm-backend/integration/osrmfiles/dotosrmdottimestamp"
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
	if flags.printSummary >= 0 {
		contents.PrintSummary(flags.printSummary)
	}

	timestampContents, err := dotosrmdottimestamp.Load(flags.filePath + ".timestamp")
	if err != nil {
		glog.Error(err)
		return
	}
	if flags.printSummary >= 0 {
		timestampContents.PrintSummary(flags.printSummary)
	}

}
