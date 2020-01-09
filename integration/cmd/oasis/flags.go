package main

import (
	"flag"
)

var flags struct {
	listenPort          int
	osrmBackendEndpoint string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8090, "Listen port.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "Backend OSRM-backend endpoint")
}
