package main

import (
	"flag"
)

var flags struct {
	listenPort           int
	osrmBackendEndpoint  string
	tnSearchEndpoint     string
	tnSearchAPIKey       string
	tnSearchAPISignature string
}

func init() {
	flag.IntVar(&flags.listenPort, "p", 8090, "Listen port.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "OSRM-backend endpoint")
	flag.StringVar(&flags.osrmBackendEndpoint, "search", "", "TN-Search-backend endpoint")
	flag.StringVar(&flags.osrmBackendEndpoint, "searchApiKey", "", "API key for TN-Search-backend")
	flag.StringVar(&flags.osrmBackendEndpoint, "searchApiSignature", "", "API Signature for  TN-Search-backend")
}
