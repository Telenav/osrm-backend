package main

import (
	"sync"
	"time"

	"github.com/Telenav/osrm-backend/integration/util/csvreader"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes"
	"github.com/Telenav/osrm-backend/integration/util/waysnodes/way2nodescsv"
	"github.com/golang/glog"
)

func pipeline(inFile string, outDBFile string, snappyCompressed bool) error {

	wayNodesChan := make(chan []waysnodes.WayNodes, 1000)
	errChan := make(chan error)
	go func() {
		errChan <- newStore(wayNodesChan, outDBFile)
	}()

	options := csvreader.DefaultOptions()
	if snappyCompressed {
		options.Compression = csvreader.CompressionTypeSnappy
	}
	l := csvreader.NewLinesAsyncReader(inFile, options)
	l.Start()

	parallel := 10 // multiple routeine to speed up parsing csv lines
	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			startTime := time.Now()
			var readCount, parseSucceedCount int
			var err error
			for {
				lines, ok := l.ReadLines()
				if !ok {
					break
				}
				readCount += len(lines)

				wayNodesSlice := make([]waysnodes.WayNodes, 0, len(lines))
				for _, line := range lines {
					wayNodes := waysnodes.WayNodes{}
					wayNodes.WayID, wayNodes.NodeIDs, err = way2nodescsv.ParseLine(line)
					if err != nil {
						glog.Warningf("Parse line %s failed, err: %v", line, err)
						continue
					}
					wayNodesSlice = append(wayNodesSlice, wayNodes)
					parseSucceedCount++
				}

				wayNodesChan <- wayNodesSlice
			}
			glog.Infof("Read/parse routine takes %f seconds for processing, read lines %d, parsed %d.", time.Now().Sub(startTime).Seconds(), readCount, parseSucceedCount)
			wg.Done()
		}()
	}
	wg.Wait()
	close(wayNodesChan)
	if err := l.Err(); err != nil {
		return err
	}

	return <-errChan
}
