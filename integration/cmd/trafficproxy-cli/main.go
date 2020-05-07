package main

import (
	"flag"

	"github.com/Telenav/osrm-backend/integration/util/appversion"

	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficdumper"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxy"
	"github.com/Telenav/osrm-backend/integration/traffic/livetraffic/trafficproxyclient"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	appversion.PrintExit()
	defer glog.Flush()

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
		h := trafficdumper.NewHandler()
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
		h := trafficdumper.NewHandler()
		h.DumpFlowResponses(trafficResp.FlowResponses)
		h.DumpIncidentResponses(trafficResp.IncidentResponses)
		return
	} else if flags.rpcMode == rpcModeStreamingDelta {

		responseChan := make(chan trafficproxy.TrafficResponse)
		waitChan := make(chan struct{})

		// async startup dumper
		go func() {
			trafficdumper.DumpStreamingDelta(responseChan)
			waitChan <- struct{}{}
		}()

		// startup streaming delta
		err := trafficproxyclient.StreamingDeltaFlowsIncidents(responseChan)
		if err != nil {
			glog.Error(err)
		}
		close(responseChan)
		<-waitChan
		return
	}

	glog.Errorf("unknown mode %s", flags.rpcMode)
}
