package main

import (
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/oasis"
	"github.com/golang/glog"
)

func main() {
	mux := http.NewServeMux()

	oasisService := oasis.New(flags.osrmBackendEndpoint)
	mux.Handle("/oasis/v1/earliest/", oasisService)

	// listen
	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
