package wayid2nodeids

import (
	"bufio"
	"os"
	"strconv"
	"strings"
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

	lineChan := make(chan string)

	// read from file
	go func() {
		scanner := bufio.NewScanner(snappy.NewReader(f))
		for scanner.Scan() {
			lineChan <- (scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			glog.Errorf("reading standard input: %v", err)
		}

		close(lineChan)
	}()

	idsChan := make(chan []int64)
	go func() {
		for {
			line, ok := <-lineChan
			if !ok {
				close(idsChan)
				break
			}
			ids := parseLine(line)
			if ids == nil {
				continue
			}

			idsChan <- ids
		}
	}()

	var nodeCount int64
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
		nodeCount += int64(len(nodeIDs))

		m.wayID2NodeIDs[wayID] = nodeIDs // store wayID->NodeID,NodeID,... mapping

		for i := range nodeIDs[:len(nodeIDs)-1] { // store Edge->wayID mapping
			edge := nodebasededge.Edge{FromNode: nodeIDs[i], ToNode: nodeIDs[i+1]}
			m.edge2WayID[edge] = wayID
		}

	}

	glog.Infof("Load wayID->nodeIDs mapping, total processing time %f seconds, ways count %d, nodes count %d.",
		time.Now().Sub(startTime).Seconds(), len(m.wayID2NodeIDs), nodeCount)

	return nil
}

func parseLine(line string) []int64 {

	elements := strings.Split(line, ",")
	if len(elements) < 3 { // at least should be one wayID and two NodeIDs
		glog.Warningf("wrong mapping line %s", line)
		return nil
	}

	wayID, err := strconv.ParseInt(elements[0], 10, 64)
	if err != nil {
		glog.Warningf("decode wayID failed from %v\n", elements)
		return nil
	}

	nodeIDs := []int64{}
	nodeElements := elements[1:]
	for _, nodeElement := range nodeElements {
		if len(nodeElement) == 0 {
			continue // the last element might be empty string
		}

		//nodeID
		nodeID, err := strconv.ParseInt(nodeElement, 10, 64)
		if err != nil {
			glog.Warningf("decode nodeID failed from %s\n", nodeElement)
			continue
		}
		nodeIDs = append(nodeIDs, nodeID)
	}
	if len(nodeIDs) < 2 {
		glog.Warningf("too less nodeIDs %v from %s", nodeIDs, line)
		return nil
	}

	return append([]int64{wayID}, nodeIDs...)
}
