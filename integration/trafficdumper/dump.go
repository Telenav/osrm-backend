package trafficdumper

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// DumpFlowResponses dump flows to file/stdout.
func (h Handler) DumpFlowResponses(flowResponses []*trafficproxy.FlowResponse) {
	if len(flowResponses) == 0 {
		return
	}

	contentChan := make(chan string)
	waitDoneChan := make(chan struct{})
	if h.writeToFile {
		go h.dumpToCSVFile("_flows", contentChan, waitDoneChan)
	}

	for _, flow := range flowResponses {
		if glog.V(3) { // verbose debug only
			glog.Infoln(flow)
		}

		if h.blockingOnly && !flow.Flow.IsBlocking() {
			continue // ignore non-blocking flow
		}

		var csvString string
		if h.humanFriendlyCSV {
			csvString = flow.Flow.HumanFriendlyCSVString()
		} else {
			csvString = flow.Flow.CSVString()
		}

		if h.writeToStdout {
			fmt.Println(csvString)
		}
		if h.writeToFile {
			contentChan <- csvString
		}
	}

	if h.writeToFile {
		close(contentChan)
		<-waitDoneChan
	}
}

// DumpIncidentResponses dump incidents to file/stdout.
func (h Handler) DumpIncidentResponses(incidentResponses []*trafficproxy.IncidentResponse) {
	if len(incidentResponses) == 0 {
		return
	}

	contentChan := make(chan string)
	waitDoneChan := make(chan struct{})
	if h.writeToFile {
		go h.dumpToCSVFile("_incidents", contentChan, waitDoneChan)
	}

	for _, incident := range incidentResponses {
		if glog.V(3) { // verbose debug only
			glog.Infoln(incident)
		}

		if h.blockingOnly && !incident.Incident.IsBlocking {
			continue // ignore non-blocking incident
		}

		var csvString string
		if h.humanFriendlyCSV {
			csvString = incident.Incident.HumanFriendlyCSVString()
		} else {
			csvString = incident.Incident.CSVString()
		}

		if h.writeToStdout {
			fmt.Println(csvString)
		}
		if h.writeToFile {
			contentChan <- csvString
		}
	}

	if h.writeToFile {
		close(contentChan)
		<-waitDoneChan
	}
}

func (h Handler) dumpToCSVFile(fileTag string, sink <-chan string, done chan<- struct{}) {
	defer close(done)
	filePath := h.dumpFileNamePrefix + fileTag + ".csv"
	startTime := time.Now()
	var dumpedLines int64
	defer func() {
		glog.Infof("dumpToCSVFile %s takes %f seconds, total lines %d.\n",
			filePath, time.Now().Sub(startTime).Seconds(), dumpedLines)
	}()

	// open file
	//outfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	outfile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	defer outfile.Close()
	defer outfile.Sync()
	if err != nil {
		glog.Error(err)
		return
	}
	glog.V(1).Infof("open output file of %s succeed.\n", filePath)

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
