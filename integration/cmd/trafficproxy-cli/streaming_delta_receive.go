package main

import (
	"fmt"
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/golang/glog"
)

func streamingDeltaReceive() {

	responseChan := make(chan proxy.TrafficResponse)

	go func() {

		dumpFileIndex := 0
		h := newResponseHandler()
		if h.writeToFile && flags.streamingDeltaSplitDumpFiles {
			h.dumpFileNamePrefix = flags.dumpFile + fmt.Sprintf("_%d", dumpFileIndex)
		}
		startTime := time.Now()
		trafficResponse := proxy.TrafficResponse{}

		for {
			resp, ok := <-responseChan

			currTime := time.Now()
			timeInterval := currTime.Sub(startTime)
			if ok && timeInterval < flags.streamingDeltaDumpInterval {
				trafficResponse.FlowResponses = append(trafficResponse.FlowResponses, resp.FlowResponses...)
				trafficResponse.IncidentResponses = append(trafficResponse.IncidentResponses, resp.IncidentResponses...)
				continue
			}

			// handle per interval
			glog.Infof("handling flows,incidents(%d,%d) from streaming delta, interval %f seconds",
				len(trafficResponse.FlowResponses), len(trafficResponse.IncidentResponses), timeInterval.Seconds())
			h.handleFlowResponses(trafficResponse.FlowResponses)
			h.handleIncidentResponses(trafficResponse.IncidentResponses)

			if !ok { // streaming delta channel no longer available, break after handling to make sure cached data processing.
				break
			}

			startTime = currTime
			if h.writeToFile && flags.streamingDeltaSplitDumpFiles {
				dumpFileIndex++
				h.dumpFileNamePrefix = flags.dumpFile + fmt.Sprintf("_%d", dumpFileIndex)
			}
		}
	}()

	// startup streaming delta
	glog.Error(trafficproxyclient.StreamingDeltaFlowsIncidents(responseChan))
}
