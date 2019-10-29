package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
)

var flags struct {
	mappingFile string
	csvFile     string
}

func init() {
	flag.StringVar(&flags.mappingFile, "m", "wayid2nodeids.csv", "OSRM way id to node ids mapping table")
	flag.StringVar(&flags.csvFile, "f", "traffic.csv", "OSRM traffic csv file")
}

const TASKNUM = 128
const CACHEDOBJECTS = 4000000

func main() {
	flag.Parse()

	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		fmt.Printf("Total processing time %f seconds\n", endTime.Sub(startTime).Seconds())
	}()

	isFlowDoneChan := make(chan bool, 1)
	wayid2speed := make(map[int64]int)
	go func() {
		trafficData, err := trafficproxyclient.GetFlowsIncidents(nil)
		if err != nil {
			log.Println(err)
			isFlowDoneChan <- false
			return
		}

		trafficData2map(*trafficData, wayid2speed)
		isFlowDoneChan <- true
	}()

	var sources [TASKNUM]chan string
	for i := range sources {
		sources[i] = make(chan string, CACHEDOBJECTS)
	}
	go loadWay2NodeidsTable(flags.mappingFile, sources)

	isFlowDone := wait4PreConditions(isFlowDoneChan)
	if isFlowDone {
		var ds dumperStatistic
		ds.Init(TASKNUM)
		dumpSpeedTable4Customize(wayid2speed, sources, flags.csvFile, &ds)
		ds.Output()
	}
}

func wait4PreConditions(flowChan <-chan bool) bool {
	var isFlowDone bool
loop:
	for {
		select {
		case f := <-flowChan:
			if !f {
				fmt.Printf("[ERROR] Communication with traffic server failed.\n")
				break loop
			} else {
				isFlowDone = true
				break loop
			}
		}
	}
	return isFlowDone
}

func trafficData2map(trafficData proxy.TrafficResponse, m map[int64]int) {
	startTime := time.Now()
	defer func() {
		log.Printf("Processing time for building traffic map takes %f seconds\n", time.Now().Sub(startTime).Seconds())
	}()

	var fwdCnt, bwdCnt uint64
	for _, flow := range trafficData.FlowResponses {
		wayid := flow.Flow.WayId
		m[wayid] = int(flow.Flow.Speed)

		if wayid > 0 {
			fwdCnt++
		} else {
			bwdCnt++
		}
	}

	//TODO: support incidents

	log.Printf("Load map[wayid] to speed with %d items, %d forward and %d backward.\n", (fwdCnt + bwdCnt), fwdCnt, bwdCnt)
}
