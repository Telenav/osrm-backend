package main

import (
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes/nodes2wayblotdb"
	"github.com/golang/glog"
)

func newDBBuilder(in <-chan []waysnodes.WayNodes, out string) (err error) {
	startTime := time.Now()

	db, err := nodes2wayblotdb.Open(out, false)
	if err != nil {
		return
	}
	defer func() {
		err = db.Close()
	}()
	var inCount, succeedCount int
	for {
		wayNodesSlice, ok := <-in
		if !ok {
			break
		}
		inCount += len(wayNodesSlice)

		written, err := writeDB(db, wayNodesSlice)
		if err != nil {
			glog.Error(err)
			break
		}
		succeedCount += written
	}

	glog.Infof("Built DB %s, in count %d, succeed count %d, takes %f seconds", out, inCount, succeedCount, time.Now().Sub(startTime).Seconds())
	return
}

func writeDB(db *nodes2wayblotdb.DB, wayNodesSlice []waysnodes.WayNodes) (int, error) {
	if db == nil {
		err := fmt.Errorf("empty db")
		glog.Fatal(err)
		return 0, err
	}

	// for best performance we possible to have
	// See more in https://github.com/Telenav/osrm-backend/issues/272#issuecomment-612877931
	const batchWriteCount = 100

	var count int
	for {
		if len(wayNodesSlice) < batchWriteCount {
			break
		}

		if err := db.BatchWrite(wayNodesSlice[:batchWriteCount]); err != nil {
			return count, fmt.Errorf("Write into db failed, err: %v", err)
		}

		count += batchWriteCount
		wayNodesSlice = wayNodesSlice[batchWriteCount:]
	}

	if len(wayNodesSlice) > 0 {
		if err := db.BatchWrite(wayNodesSlice); err != nil {
			return count, fmt.Errorf("Write into db failed, err: %v", err)
		}
		count += len(wayNodesSlice)
	}

	return count, nil
}
