package main

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes/nodes2wayblotdb"
	"github.com/golang/glog"
)

func newStore(in <-chan *waysnodes.WayNodes, out string) (err error) {
	startTime := time.Now()

	db, err := nodes2wayblotdb.Open(out, false)
	if err != nil {
		return
	}
	defer func() {
		err = db.Close()
	}()

	batchCountPerWrite := 100
	wayNodesCache := make([]waysnodes.WayNodes, batchCountPerWrite)
	var inCount, succeedCount int
	for {
		wayNodes, ok := <-in
		if !ok {
			break
		}
		inCount++
		wayNodesCache = append(wayNodesCache, *wayNodes)

		if len(wayNodesCache) < batchCountPerWrite {
			continue
		}

		if err := db.BatchWrite(wayNodesCache); err != nil {
			glog.Errorf("Write into db failed, err: %v", err)
			break
		}
		succeedCount += len(wayNodesCache)
		wayNodesCache = wayNodesCache[:0]
	}

	glog.V(1).Infof("Built DB %s, in count %d, succeed count %d, takes %f seconds", out, inCount, succeedCount, time.Now().Sub(startTime).Seconds())
	return
}
