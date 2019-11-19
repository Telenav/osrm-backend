package wayid2nodeids

import (
	"bufio"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/golang/snappy"
)

func (m *Mapping) load() error {
	startTime := time.Now()

	f, err := os.Open(m.mappingFile)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(2).Infof("open wayid2nodeids mapping file %s succeed.\n", m.mappingFile)

	//// reversely start tasks

	// start task: store IDs to map
	idsChan := make(chan []int64, 100) // use channel with cache here to significantly improve performance
	waitStoreTaskChan := make(chan struct{})
	go m.storeTask(idsChan, waitStoreTaskChan)

	// start task: parse line string to ID slice
	const parseLineTaskCount = 8
	lineChan := make(chan string, parseLineTaskCount)
	waitParseTaskChan := make(chan struct{}, parseLineTaskCount)
	for i := 0; i < parseLineTaskCount; i++ {
		go parseLineTask(lineChan, idsChan, waitParseTaskChan)
	}

	// start task: read from file
	waitReadDone := make(chan error)
	go m.readTask(f, lineChan, waitReadDone)

	// wait done
	readErr := <-waitReadDone
	if readErr != nil {
		glog.Warning(readErr)
	}
	close(lineChan)
	for i := 0; i < parseLineTaskCount; i++ {
		<-waitParseTaskChan
	}
	close(idsChan)
	<-waitStoreTaskChan

	glog.Infof("Load wayID->nodeIDs mapping, total processing time %f seconds, loaded ways %d.",
		time.Now().Sub(startTime).Seconds(), len(m.wayID2NodeIDs))

	return readErr
}

func (m *Mapping) readTask(f *os.File, lineChan chan<- string, done chan<- error) {
	if f == nil {
		glog.Fatalf("file %v invalid", f)
	}
	var lineCount int64

	// start task: read from file
	scanner := bufio.NewScanner(snappy.NewReader(f))
	for scanner.Scan() {
		lineChan <- (scanner.Text())
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		glog.Error(err)
		done <- err
	}
	done <- nil
	glog.V(2).Infof("read file task done, total read line count %d", lineCount)
}

func (m *Mapping) storeTask(idsChan <-chan []int64, done chan<- struct{}) {
	var nodesCountWithinDuplicate, possibleEdgesCount int64
	for {
		ids, ok := <-idsChan
		if !ok {
			break
		}

		if len(ids) < 3 {
			glog.Errorf("expect at least 3 ids(wayID,nodeID,nodeID) but not enough: %v", ids)
			continue
		}

		wayID := ids[0]
		nodeIDs := ids[1:]
		nodesCountWithinDuplicate += int64(len(nodeIDs))
		possibleEdgesCount += int64(len(nodeIDs) - 1)

		m.wayID2NodeIDs[wayID] = nodeIDs // store wayID->NodeID,NodeID,... mapping

	}
	done <- struct{}{}
	glog.V(2).Infof("store into map task done, total ways %d, nodes(within duplicate) %d, edges(possible) %d", len(m.wayID2NodeIDs), nodesCountWithinDuplicate, possibleEdgesCount)
}

func parseLineTask(lineChan <-chan string, result chan<- []int64, done chan<- struct{}) {
	for {
		line, ok := <-lineChan
		if !ok {
			break
		}

		ids := parseLine(line)
		if ids == nil {
			continue
		}

		result <- ids

	}
	done <- struct{}{}
}
