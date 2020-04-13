package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/service/oasis"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	mux := http.NewServeMux()

	oasisService, err := oasis.New(flags.osrmBackendEndpoint, flags.finderType, flags.tnSearchEndpoint,
		flags.tnSearchAPIKey, flags.tnSearchAPISignature, flags.localDataPath)
	if err != nil {
		glog.Errorf("Failed to create oasis handler due to err %+v.\n", err)
		return
	}

	mux.Handle("/oasis/v1/earliest/", oasisService)

	// listen
	listening := ":" + strconv.Itoa(flags.listenPort)
	glog.Infof("Listening on %s", listening)
	glog.Fatal(http.ListenAndServe(listening, mux))
}
