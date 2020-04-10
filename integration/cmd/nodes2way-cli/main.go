package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes/nodes2wayblotdb"
	"github.com/golang/glog"
)

// output logs to stderr without timestamp
var cliLog = log.New(os.Stderr, "", 0)

func main() {
	flag.Parse()

	wayIDs, err := query(flags.db, flags.nodeIDs)
	if err != nil {
		cliLog.Println(err)
		os.Exit(1)
		return
	}
	fmt.Println(wayIDs)
}

func query(dbFile string, nodeIDs []int64) ([]int64, error) {

	db, err := nodes2wayblotdb.Open(dbFile, true)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := time.Now()

	wayIDs, err := db.QueryWays(nodeIDs)
	if err != nil {
		return nil, err
	}

	glog.Infof("Querying takes %f seconds", time.Now().Sub(startTime).Seconds()) // easy to measure querying time cost

	return wayIDs, nil
}
