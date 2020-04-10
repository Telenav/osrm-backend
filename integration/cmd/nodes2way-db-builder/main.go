package main

import (
	"flag"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	pipeline()
}

func pipeline() {
	startTime := time.Now()

	wayNodesChan := make(chan *waysnodes.WayNodes)
	errChan := make(chan error)

	go func() {
		errChan <- newStore(wayNodesChan, flags.out)
	}()
	go func() {
		errChan <- newReader(flags.in, flags.snappyCompressed, wayNodesChan)
		close(wayNodesChan)
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			glog.Error(err)
			os.Exit(1)
			return
		}
	}

	glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Now().Sub(startTime).Seconds())
}
