package main

import (
	"flag"
)

var flags struct {
	entityEndpoint string
	apiKey         string
	apiSignature   string
}

func init() {
	flag.StringVar(&flags.entityEndpoint, "entity", "", "Backend entity service endpoint")
	flag.StringVar(&flags.apiKey, "apikey", "", "API key for entity service")
	flag.StringVar(&flags.apiSignature, "apisignature", "", "API signature for entity service")
}
