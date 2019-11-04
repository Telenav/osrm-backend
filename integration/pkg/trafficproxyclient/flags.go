package trafficproxyclient

import "flag"

// Flags represent parameters to get or streaming teleanv traffic.
type Flags struct {
	Port            int
	IP              string
	Region          string
	TrafficProvider string
	MapProvider     string
	flow            bool
	incident        bool
}

var flags Flags

func init() {
	flag.IntVar(&flags.Port, "p", 10086, "target traffic proxy port")
	flag.StringVar(&flags.IP, "c", "127.0.0.1", "target traffic proxy ip address")
	flag.StringVar(&flags.Region, "region", "na", "region")
	flag.StringVar(&flags.TrafficProvider, "traffic", "", "traffic data provider")
	flag.StringVar(&flags.MapProvider, "map", "", "map data provider")
	flag.BoolVar(&flags.flow, "flow", true, "Enable traffic flow.")
	flag.BoolVar(&flags.incident, "incident", true, "Enable traffic incident.")
}
