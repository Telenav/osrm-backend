package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
	"github.com/golang/glog"
)

type responseHandler struct {
	blockingOnly       bool
	writeToStdout      bool
	writeToFile        bool
	dumpFileNamePrefix string
}

func newResponseHandler() responseHandler {
	h := responseHandler{}
	h.blockingOnly = flags.blockingOnly
	h.writeToStdout = flags.stdout
	if len(flags.dumpFile) > 0 {
		h.writeToFile = true
		h.dumpFileNamePrefix = flags.dumpFile
	} else {
		h.writeToFile = false
	}
	return h
}

func (r responseHandler) handleFlowResponses(flowResponses []*proxy.FlowResponse) {

	contentChan := make(chan string)
	waitDoneChan := make(chan struct{})
	if r.writeToFile {
		go r.dumpToCSVFile("_flows", contentChan, waitDoneChan)
	}

	for _, flow := range flowResponses {
		if glog.V(3) { // verbose debug only
			glog.Infoln(flow)
		}

		if r.blockingOnly && !flow.Flow.IsBlocking() {
			continue // ignore non-blocking flow
		}

		if r.writeToStdout {
			fmt.Println(flow.Flow.CSVString())
		}
		if r.writeToFile {
			contentChan <- flow.Flow.CSVString()
		}
	}

	if r.writeToFile {
		close(contentChan)
		<-waitDoneChan
	}
}

func (r responseHandler) handleIncidentResponses(incidentResponses []*proxy.IncidentResponse) {

	contentChan := make(chan string)
	waitDoneChan := make(chan struct{})
	if r.writeToFile {
		go r.dumpToCSVFile("_incidents", contentChan, waitDoneChan)
	}

	for _, incident := range incidentResponses {
		if glog.V(3) { // verbose debug only
			glog.Infoln(incident)
		}

		if r.blockingOnly && !incident.Incident.IsBlocking {
			continue // ignore non-blocking incident
		}

		if r.writeToStdout {
			fmt.Println(incident.Incident.CSVString())
		}
		if r.writeToFile {
			contentChan <- incident.Incident.CSVString()
		}
	}

	if r.writeToFile {
		close(contentChan)
		<-waitDoneChan
	}
}

func (r responseHandler) dumpToCSVFile(fileTag string, sink <-chan string, done chan<- struct{}) {
	defer close(done)
	filePath := r.dumpFileNamePrefix + fileTag + ".csv"
	startTime := time.Now()
	var dumpedLines int64
	defer func() {
		glog.Infof("dumpToCSVFile %s takes %f seconds, total lines %d.\n",
			filePath, time.Now().Sub(startTime).Seconds(), dumpedLines)
	}()

	// open file
	outfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	defer outfile.Close()
	defer outfile.Sync()
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Infof("open output file of %s succeed.\n", filePath)

	// write contents
	w := bufio.NewWriter(outfile)
	defer w.Flush()
	for {
		str, ok := <-sink
		if !ok {
			break // gracefully done
		}

		_, err := w.WriteString(str + "\n")
		if err != nil {
			glog.Error(err)
			return
		}
		dumpedLines++
	}
}
