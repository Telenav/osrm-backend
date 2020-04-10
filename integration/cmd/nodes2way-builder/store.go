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

	var inCount, succeedCount int
	for {
		wayNodes, ok := <-in
		if !ok {
			break
		}
		inCount++

		if err := db.Write(wayNodes.WayID, wayNodes.NodeIDs); err != nil {
			if glog.V(3) { // avoid affect performance by verbose log
				glog.Infof("Update %+v into db failed, err: %v", wayNodes, err)
			}
			continue
		}
		succeedCount++
	}

	glog.V(1).Infof("Built DB %s, in count %d, succeed count %d, takes %f seconds", out, inCount, succeedCount, time.Now().Sub(startTime).Seconds())
	return
}
