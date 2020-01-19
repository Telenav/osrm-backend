package trafficdumper

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// DumpStreamingDelta dumps traffic response from streaming delta channel.
func DumpStreamingDelta(responseChan <-chan trafficproxy.TrafficResponse) {

	h := NewHandler()
	startTime := time.Now()
	trafficResponse := trafficproxy.TrafficResponse{}

	for {
		resp, ok := <-responseChan

		currTime := time.Now()
		timeInterval := currTime.Sub(startTime)
		if ok && timeInterval < h.streamingDeltaDumpInterval {
			trafficResponse.FlowResponses = append(trafficResponse.FlowResponses, resp.FlowResponses...)
			trafficResponse.IncidentResponses = append(trafficResponse.IncidentResponses, resp.IncidentResponses...)
			continue
		}

		// handle per interval
		glog.Infof("handling flows,incidents(%d,%d) from streaming delta, interval %f seconds",
			len(trafficResponse.FlowResponses), len(trafficResponse.IncidentResponses), timeInterval.Seconds())
		if h.writeToFile && h.streamingDeltaSplitDumpFiles {
			h.updateDumpFileNamePrefix()
		}
		h.DumpFlowResponses(trafficResponse.FlowResponses)
		h.DumpIncidentResponses(trafficResponse.IncidentResponses)

		if !ok { // streaming delta channel no longer available, break after handling to make sure cached data processing.
			break
		}

		startTime = currTime
	}
}

// updateDumpFileNamePrefix updates prefix for next dump splited files.
func (h *Handler) updateDumpFileNamePrefix() {
	h.dumpFileNamePrefix = flags.dumpFile + fmt.Sprintf("_%d", h.dumpFileSplitIndex)
	h.dumpFileSplitIndex++
}
