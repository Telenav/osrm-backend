package main

import (
	"flag"
	"fmt"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxyclient"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	if !flags.streamingDelta {
		trafficResp, err := trafficproxyclient.GetFlowsIncidents(flags.wayIDsFlag.wayIDs)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infof("total received traffic flows,incidents(%d,%d)",
			len(trafficResp.FlowResponses), len(trafficResp.IncidentResponses))

		for _, flow := range trafficResp.FlowResponses {
			if glog.V(3) { // verbose debug only
				glog.Infoln(flow)
			}
			fmt.Println(flow.Flow.CSVString())
		}
		for _, incident := range trafficResp.IncidentResponses {
			if glog.V(3) { // verbose debug only
				glog.Infoln(incident)
			}
			fmt.Println(incident.Incident.CSVString())
		}

		return
	}

}
