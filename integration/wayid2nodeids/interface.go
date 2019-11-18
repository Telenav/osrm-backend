package wayid2nodeids

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/golang/snappy"
)

// Mapping handles 'wayID->NodeID,NodeID,NodeID,...' mapping.
type Mapping struct {
	mappingFile string
	m           map[int64][]int64
}

// NewMappingFrom creates a new Mapping object for 'wayID->NodeID,NodeID,NodeID,...' mapping.
// Currently it only supports mapping file compressed by snappy, e.g. 'wayid2nodeids.csv.snappy'.
func NewMappingFrom(mappingFilePath string) Mapping {
	m := Mapping{
		mappingFilePath, map[int64][]int64{},
	}
	return m
}

// Load loads data from file to map in memory, it will returns until the whole load process done.
func (m *Mapping) Load() error {
	startTime := time.Now()

	f, err := os.Open(m.mappingFile)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(2).Infof("open wayid2nodeids mapping file %s succeed.\n", m.mappingFile)

	lineChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(snappy.NewReader(f))
		for scanner.Scan() {
			lineChan <- (scanner.Text())
		}
		close(lineChan)
	}()

	var wayCount, nodeCount int64
	for {
		line, ok := <-lineChan
		if !ok {
			break
		}

		elements := strings.Split(line, ",")
		if len(elements) < 3 { // at least should be one wayID and two NodeIDs
			glog.Warningf("wrong mapping line %s", line)
			continue
		}

		wayID, err := strconv.ParseInt(elements[0], 10, 64)
		if err != nil {
			glog.Warningf("decode wayID failed from %v\n", elements)
			continue
		}
		m.m[wayID] = []int64{}
		wayCount++

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
			m.m[wayID] = append(m.m[wayID], nodeID)
			nodeCount++
		}

	}

	glog.Infof("Load wayID->nodeIDs mapping, total processing time %f seconds, ways count %d, nodes count %d.",
		time.Now().Sub(startTime).Seconds(), wayCount, nodeCount)
	return nil
}
