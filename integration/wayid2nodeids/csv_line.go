package wayid2nodeids

import (
	"strconv"
	"strings"

	"github.com/golang/glog"
)

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
