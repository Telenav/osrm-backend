package main

import (
	"flag"
	"strings"

	"github.com/Telenav/osrm-backend/integration/service/ranking/trafficapplyingmodel/modelfactory"
)

var flags struct {
	listenPort int

	nodes2WayDB string

	liveTraffic                     bool // whether enable live traffic or not
	historicalSpeed                 bool // whether enable historical speed or not
	historicalSpeedDailyPatternFile string
	historicalSpeedWaysMappingFile  string

	trafficApplyingModel string

	osrmBackendEndpoint string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8080, "Listen port.")

	flag.StringVar(&flags.nodes2WayDB, "nodes2way", "nodes2way.db", "BoltDB for querying wayIDs from nodeIDs.")

	flag.BoolVar(&flags.liveTraffic, "live-traffic", false, "Enable live traffic. ")
	flag.BoolVar(&flags.historicalSpeed, "historical-speed", false, "Enable historical speed. The historical speed related files won't be loaded if disabled.")
	flag.StringVar(&flags.historicalSpeedDailyPatternFile, "historical-speed-dailypattern", "", "Historical speed daily patterns csv file.")
	flag.StringVar(&flags.historicalSpeedWaysMappingFile, "historical-speed-waysmapping", "", "Historical speed wayIDs to daily patterns mapping csv file. Pass in multiple files separated by ','.")

	flag.StringVar(&flags.trafficApplyingModel, "traffic-applying-model", modelfactory.DefaultModelName(), "Model to apply traffic on route, options: "+strings.Join(modelfactory.AvailableModelNames(), ","))

	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "Backend OSRM-backend endpoint")
}
