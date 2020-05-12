package main

import (
	"flag"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/service/oasis/connectivitymap"
	"github.com/Telenav/osrm-backend/integration/service/oasis/osrmconnector"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/ranker"
	"github.com/Telenav/osrm-backend/integration/service/oasis/spatialindexer/s2indexer"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	startTime := time.Now()

	if flags.inputFile == "" || flags.outputFolder == "" {
		glog.Fatal("Empty string for inputFile or outputFolder, please check your input.\n")
	}

	var rankerStrategy spatialindexer.Ranker
	if flags.osrmBackendEndpoint == "" {
		glog.Warning("No information about OSRM Endpoint, can only init ranker with great circle distance.")
		rankerStrategy = ranker.CreateRanker(ranker.SimpleRanker, nil)
	} else {
		rankerStrategy = ranker.CreateRanker(ranker.OSRMBasedRanker,
			osrmconnector.NewOSRMConnector(flags.osrmBackendEndpoint))
	}

	indexer := s2indexer.NewS2Indexer().Build(flags.inputFile)
	if indexer == nil {
		glog.Fatalf("Failed to build indexer, stop %s\n", os.Args[0])
	}
	indexer.Dump(flags.outputFolder)

	connectivitymap.New(flags.maxRange).
		Build(indexer, indexer, rankerStrategy, flags.numberOfWorkers).
		Dump(flags.outputFolder)

	glog.Infof("%s totally takes %f seconds for processing.", os.Args[0], time.Since(startTime).Seconds())
}
