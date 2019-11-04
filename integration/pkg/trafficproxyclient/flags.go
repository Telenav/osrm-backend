package trafficproxyclient

import (
	"flag"
	"time"
)

var flags struct {
	port            int
	ip              string
	region          string
	trafficProvider string
	mapProvider     string
	flow            bool
	incident        bool
	rpcGetTimeout   time.Duration
}

func init() {
	flag.IntVar(&flags.port, "p", 10086, "target traffic proxy port")
	flag.StringVar(&flags.ip, "c", "127.0.0.1", "target traffic proxy ip address")
	flag.StringVar(&flags.region, "region", "na", "region")
	flag.StringVar(&flags.trafficProvider, "traffic", "", "traffic data provider")
	flag.StringVar(&flags.mapProvider, "map", "", "map data provider")
	flag.BoolVar(&flags.flow, "flow", true, "Enable traffic flow.")
	flag.BoolVar(&flags.incident, "incident", true, "Enable traffic incident.")
	flag.DurationVar(&flags.rpcGetTimeout, "gettimeout", 60*time.Second, "Timeout for getting via RPC.")
}
