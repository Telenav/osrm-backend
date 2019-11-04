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

		for _, flow := range trafficResp.FlowResponses {
			fmt.Println(flow)
		}
		for _, incident := range trafficResp.IncidentResponses {
			fmt.Println(incident)
		}

		return
	}

}
