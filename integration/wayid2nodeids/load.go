package wayid2nodeids

import (
	"bufio"
	"os"
	"time"

	"github.com/Telenav/osrm-backend/integration/nodebasededge"
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
	const storeTaskCount = 3
	idsChan := []chan []int64{}
	for i := 0; i < storeTaskCount; i++ {
		idsChan = append(idsChan, make(chan []int64, 100))
	}
	waitStoreTaskChan := make(chan struct{}, storeTaskCount)
	go m.storeWayID2NodeIDs(idsChan[0], waitStoreTaskChan)
	go m.storeEdge2WayID(idsChan[1], waitStoreTaskChan)
	go m.storeNodeIDs(idsChan[2], waitStoreTaskChan)

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
	for i := range idsChan {
		close(idsChan[i])
	}
	for i := 0; i < storeTaskCount; i++ {
		<-waitStoreTaskChan
	}

	glog.Infof("Load wayID->nodeIDs mapping, total processing time %f seconds, loaded ways %d, nodes %d, edges %d.",
		time.Now().Sub(startTime).Seconds(), len(m.wayID2NodeIDs), len(m.edge2WayID), len(m.nodeIDs))

	return readErr
}

func (m *Mapping) readTask(f *os.File, lineChan chan<- string, done chan<- error) {
	if f == nil {
		glog.Fatalf("file %v invalid", f)
	}

	// start task: read from file
	scanner := bufio.NewScanner(snappy.NewReader(f))
	for scanner.Scan() {
		lineChan <- (scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		glog.Error(err)
		done <- err
	}
	done <- nil
}

func (m *Mapping) storeWayID2NodeIDs(idsChan <-chan []int64, done chan<- struct{}) {
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

		m.wayID2NodeIDs[wayID] = nodeIDs // store wayID->NodeID,NodeID,... mapping

	}
	done <- struct{}{}
}

func (m *Mapping) storeEdge2WayID(idsChan <-chan []int64, done chan<- struct{}) {
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

		for i := range nodeIDs[:len(nodeIDs)-1] { // store Edge->wayID mapping
			edge := nodebasededge.Edge{FromNode: nodeIDs[i], ToNode: nodeIDs[i+1]}
			m.edge2WayID[edge] = wayID
		}
	}
	done <- struct{}{}
}

func (m *Mapping) storeNodeIDs(idsChan <-chan []int64, done chan<- struct{}) {
	for {
		ids, ok := <-idsChan
		if !ok {
			break
		}

		nodeIDs := ids[1:]

		for i := range nodeIDs {
			m.nodeIDs[nodeIDs[i]] = struct{}{}
		}
	}
	done <- struct{}{}
}

func parseLineTask(lineChan <-chan string, result []chan []int64, done chan<- struct{}) {
	for {
		line, ok := <-lineChan
		if !ok {
			break
		}

		ids := parseLine(line)
		if ids == nil {
			continue
		}

		for i := range result {
			result[i] <- ids
		}
	}
	done <- struct{}{}
}
