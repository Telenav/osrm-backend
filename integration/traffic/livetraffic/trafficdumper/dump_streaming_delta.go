package trafficdumper

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/golang/glog"
)

// DumpStreamingDelta dumps traffic response from streaming delta channel.
func DumpStreamingDelta(responseChan <-chan trafficproxy.TrafficResponse) {

	h := NewHandler()
	startTime := time.Now()
	trafficResponse := trafficproxy.TrafficResponse{}
	var totalFlowsCount, totalIncidentsCount uint64

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
		totalFlowsCount += uint64(len(trafficResponse.FlowResponses))
		totalIncidentsCount += uint64(len(trafficResponse.IncidentResponses))
		glog.Infof("handling flows,incidents(%d,%d) from streaming delta, interval %f seconds. Totally handled flows,incidents(%d,%d) so far.",
			len(trafficResponse.FlowResponses), len(trafficResponse.IncidentResponses), timeInterval.Seconds(), totalFlowsCount, totalIncidentsCount)
		if h.writeToFile && h.streamingDeltaSplitDumpFiles {
			h.updateDumpFileNamePrefix()
		}
		h.DumpFlowResponses(trafficResponse.FlowResponses)
		h.DumpIncidentResponses(trafficResponse.IncidentResponses)
		trafficResponse = trafficproxy.TrafficResponse{} // clean up

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
