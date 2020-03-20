package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/golang/glog"
	"github.com/qedus/osmpbf"
)

func newPBFParser(filePath string, out chan<- *osmpbf.Way) error {
	startTime := time.Now()

	pbfFile, err := os.Open(filePath)
	defer pbfFile.Close()
	if err != nil {
		return fmt.Errorf("Open pbf file of %s failed, err: %v", filePath, err)
	}
	glog.V(1).Infof("Open pbf file of %s succeed.\n", filePath)

	d := osmpbf.NewDecoder(pbfFile)
	d.SetBufferSize(osmpbf.MaxBlobSize)
	if err := d.Start(runtime.GOMAXPROCS(-1)); err != nil {
		return fmt.Errorf("start pbf decoder failed, err: %v", err)
	}

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nc++
			case *osmpbf.Way:
				// Process Way v.
				wc++
				out <- v
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				return fmt.Errorf("unknown type %T", v)
			}
		}
	}

	glog.Infof("Parsed pbf ways %d, nodes %d, relations: %d, takes %f seconds", wc, nc, rc, time.Now().Sub(startTime).Seconds())
	return nil
}
