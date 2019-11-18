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

	const parseLineTaskCount = 8
	lineChan := make(chan string, parseLineTaskCount)

	// start task: read from file
	go func() {
		scanner := bufio.NewScanner(snappy.NewReader(f))
		for scanner.Scan() {
			lineChan <- (scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			glog.Error(err)
		}

		close(lineChan)
	}()

	// start task: store IDs to map
	const storeTaskCount = 2
	idsChan := []chan []int64{}
	for i := 0; i < storeTaskCount; i++ {
		idsChan = append(idsChan, make(chan []int64))
	}
	waitStoreTaskChan := make(chan struct{}, storeTaskCount)
	go m.storeWayID2NodeIDs(idsChan[0], waitStoreTaskChan)
	go m.storeEdge2WayID(idsChan[1], waitStoreTaskChan)

	// start task: parse line string to ID slice
	waitParseTaskChan := make(chan struct{}, parseLineTaskCount)
	for i := 0; i < parseLineTaskCount; i++ {
		go parseLineTask(lineChan, idsChan, waitParseTaskChan)
	}

	// wait done
	for i := 0; i < parseLineTaskCount; i++ {
		<-waitParseTaskChan
	}
	for i := range idsChan {
		close(idsChan[i])
	}
	for i := 0; i < storeTaskCount; i++ {
		<-waitStoreTaskChan
	}

	glog.Infof("Load wayID->nodeIDs mapping, total processing time %f seconds, ways count %d.",
		time.Now().Sub(startTime).Seconds(), len(m.wayID2NodeIDs))

	return nil
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
