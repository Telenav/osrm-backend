package main

import (
	"flag"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/Telenav/osrm-backend/integration/trafficdumper"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	if flags.rpcMode == rpcModeGetWays {
		if len(flags.wayIDs) == 0 {
			glog.Error("please provide wayIDs for 'getways' mode by '-ways xxx', e.g. '-ways 829733412,-104489539'")
			return
		}

		// get traffic data
		trafficResp, err := trafficproxyclient.GetFlowsIncidents(flags.wayIDs)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infof("total received traffic flows,incidents(%d,%d)",
			len(trafficResp.FlowResponses), len(trafficResp.IncidentResponses))

		// dump traffic data
		h := trafficdumper.New()
		h.DumpFlowResponses(trafficResp.FlowResponses)
		h.DumpIncidentResponses(trafficResp.IncidentResponses)
		return
	} else if flags.rpcMode == rpcModeGetAll {

		// get traffic data
		trafficResp, err := trafficproxyclient.GetFlowsIncidents(nil)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infof("total received traffic flows,incidents(%d,%d)",
			len(trafficResp.FlowResponses), len(trafficResp.IncidentResponses))

		// dump traffic data
		h := trafficdumper.New()
		h.DumpFlowResponses(trafficResp.FlowResponses)
		h.DumpIncidentResponses(trafficResp.IncidentResponses)
		return
	} else if flags.rpcMode == rpcModeStreamingDelta {

		responseChan := make(chan proxy.TrafficResponse)

		// async startup dumper
		go func() {
			trafficdumper.DumpStreamingDelta(responseChan)
		}()

		// startup streaming delta
		err := trafficproxyclient.StreamingDeltaFlowsIncidents(responseChan)
		if err != nil {
			glog.Error(err)
		}
		return
	}

	glog.Errorf("unknown mode %s", flags.rpcMode)
}
