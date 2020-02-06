package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/oasis"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	mux := http.NewServeMux()

	oasisService := oasis.New(flags.osrmBackendEndpoint)
	mux.Handle("/oasis/v1/earliest/", oasisService)

	// listen
	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
