package rankingservice

import "flag"

var flags struct {
	alternatives uint
}

func init() {
	flag.UintVar(&flags.alternatives, "alternatives", 3, "Query the number of routes from backend OSRM.")
}
