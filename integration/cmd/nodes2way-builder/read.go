package main

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/golang/snappy"

	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes/way2nodescsv"
	"github.com/golang/glog"
)

func newReader(in string, snappyCompressed bool, out chan<- *waysnodes.WayNodes) error {
	startTime := time.Now()

	f, err := os.Open(in)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", in)

	var r *csv.Reader
	if snappyCompressed {
		r = csv.NewReader(snappy.NewReader(f))
	} else {
		r = csv.NewReader(f)
	}
	r.ReuseRecord = true
	r.FieldsPerRecord = -1 // disable fields count check

	var total, succeed int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		total++

		wayNodes := waysnodes.WayNodes{}
		wayNodes.WayID, wayNodes.NodeIDs, err = way2nodescsv.ParseRecord(record)
		if err != nil {
			glog.Warningf("Parse record %v failed, err: %v", record, err)
			continue
		}
		out <- &wayNodes
		succeed++
	}

	glog.V(1).Infof("Read wayID,nodeID,nodeID,... from file %s, total count %d, succeed parsed %d, takes %f seconds", in, total, succeed, time.Now().Sub(startTime).Seconds())
	return nil
}
