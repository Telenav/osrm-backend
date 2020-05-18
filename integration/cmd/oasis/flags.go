package main

import (
	"flag"
)

var flags struct {
	listenPort           int
	osrmBackendEndpoint  string
	finderType           string
	tnSearchEndpoint     string
	tnSearchAPIKey       string
	tnSearchAPISignature string
	localDataPath        string
	cpuProfileFile       string
}

func init() {

	flag.IntVar(&flags.listenPort, "p", 8090, "Listen port.")
	flag.StringVar(&flags.osrmBackendEndpoint, "osrm", "", "OSRM-backend endpoint")
	flag.StringVar(&flags.finderType, "finder", "", "Specify search finder to search for nearby charge stations for given location, use CloudFinder or LocalFinder")
	flag.StringVar(&flags.tnSearchEndpoint, "search", "", "TN-Search-backend endpoint")
	flag.StringVar(&flags.tnSearchAPIKey, "searchApiKey", "", "API key for TN-Search-backend")
	flag.StringVar(&flags.tnSearchAPISignature, "searchApiSignature", "", "API Signature for TN-Search-backend")
	flag.StringVar(&flags.localDataPath, "datapath", "", "Local data path for index data")
	flag.StringVar(&flags.cpuProfileFile, "cpuprofile", "", "write cpu profile to `file`")
}
