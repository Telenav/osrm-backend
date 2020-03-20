package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/qedus/osmpbf"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	if flags.pbfSource != pbfSourceUniDB {
		glog.Errorf("Only support %s at the moment.", pbfSourceUniDB)
		os.Exit(1)
		return
	}

	startTime := time.Now()

	osmWaysChan := make(chan *osmpbf.Way)
	wayTimezoneInfoChan := make(chan *wayTimezoneInfo)
	errChan := make(chan error)

	go func() {
		glog.V(1).Infof("ways mapping updater routine start.")
		errChan <- newWaysMappingUpdater(strings.Split(flags.inWaysMappingFile, ","), flags.outWaysMappingFile, flags.withCSVHeader, wayTimezoneInfoChan)
		glog.V(1).Info("ways mapping updater routine exited.")
	}()
	go func() {
		glog.V(1).Infof("pbf parser routine start.")
		errChan <- newPBFParser(flags.pbf, osmWaysChan)
		close(osmWaysChan)
		glog.V(1).Info("pbf parser routine exited.")
	}()
	go func() {
		glog.V(1).Infof("timezone builder for unidb routine start.")
		newTimezoneBuilderForUniDB(osmWaysChan, wayTimezoneInfoChan)
		close(wayTimezoneInfoChan)
		errChan <- nil
		glog.V(1).Info("timezone builder for unidb routine exited.")
	}()

	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			glog.Error(err)
			os.Exit(1)
			return
		}
	}

	glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Now().Sub(startTime).Seconds())
}
