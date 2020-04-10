package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes/nodes2wayblotdb"
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

	wayIDs, err := db.QueryWays(nodeIDs)
	if err != nil {
		return nil, err
	}

	return wayIDs, nil
}
