package main

import (
	"flag"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/appversion"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	appversion.PrintExit()
	defer glog.Flush()

	startTime := time.Now()

	err := pipeline(flags.in, flags.out, flags.snappyCompressed)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Now().Sub(startTime).Seconds())
}
