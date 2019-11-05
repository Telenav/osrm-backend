package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	if flags.rpcMode == rpcModeGetWays {
		if len(flags.wayIDs) == 0 {
			glog.Error("please provide wayIDs for 'getways' mode by '-ways xxx', e.g. '-ways 829733412,-104489539'")
			return
		}

		trafficResp, err := trafficproxyclient.GetFlowsIncidents(flags.wayIDs)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infof("total received traffic flows,incidents(%d,%d)",
			len(trafficResp.FlowResponses), len(trafficResp.IncidentResponses))

		h := newResponseHandler()
		h.handleFlowResponses(trafficResp.FlowResponses)
		h.handleIncidentResponses(trafficResp.IncidentResponses)
		return
	} else if flags.rpcMode == rpcModeGetAll {

		trafficResp, err := trafficproxyclient.GetFlowsIncidents(nil)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infof("total received traffic flows,incidents(%d,%d)",
			len(trafficResp.FlowResponses), len(trafficResp.IncidentResponses))

		h := newResponseHandler()
		h.handleFlowResponses(trafficResp.FlowResponses)
		h.handleIncidentResponses(trafficResp.IncidentResponses)
		return
	} else if flags.rpcMode == rpcModeStreamingDelta {

		streamingDeltaReceive()
		return
	}

	glog.Errorf("unknown mode %s", flags.rpcMode)
}
