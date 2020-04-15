package main

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	startTime := time.Now()
	defer glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Now().Sub(startTime).Seconds())

	err := pipeline(flags.in, flags.out, flags.snappyCompressed)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
