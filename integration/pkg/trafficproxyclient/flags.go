package trafficproxyclient

import (
	"flag"
	"time"
)

var flags struct {
	target                string
	region                string
	trafficProvider       string
	mapProvider           string
	flow                  bool
	incident              bool
	rpcGetTimeout         time.Duration
	streamingDeltaMaxSize int
	streamingDeltaMaxTime time.Duration
}

func init() {
	flag.StringVar(&flags.target, "c", "127.0.0.1:10086", "target traffic proxy endpoint, format 'ip[:port]'. Default port 10086 will be used if only ip input, i.e '127.0.0.1' is same with '127.0.0.1:10086'.")
	flag.StringVar(&flags.region, "region", "na", "region")
	flag.StringVar(&flags.trafficProvider, "traffic", "", "traffic data provider")
	flag.StringVar(&flags.mapProvider, "map", "", "map data provider")
	flag.BoolVar(&flags.flow, "flow", true, "Enable traffic flow.")
	flag.BoolVar(&flags.incident, "incident", true, "Enable traffic incident.")
	flag.DurationVar(&flags.rpcGetTimeout, "gettimeout", 60*time.Second, "Timeout for getting via RPC.")
	flag.IntVar(&flags.streamingDeltaMaxSize, "stream-maxsize", 10000, "Max flows count per streaming delta send.")
	flag.DurationVar(&flags.streamingDeltaMaxTime, "stream-maxtime", time.Second, "Max time interval per streaming delta send.")
}
