// Package way2nodescsv provides utilities to handle way2nodes mapping that stores by csv format/file.
// Each record format is `wayID->nodeID,nodeID,nodeID,...`.
package way2nodescsv

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRecord parses way2nodes csv record.
func ParseRecord(record []string) (int64, []int64, error) {
	if len(record) < 3 { // at least should be one wayID and two NodeIDs
		return 0, nil, fmt.Errorf("invalid record %v", record)
	}

	wayID, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("decode wayID failed from %v", record)
	}

	nodeIDs := []int64{}
	nodeElements := record[1:]
	for _, nodeElement := range nodeElements {
		if len(nodeElement) == 0 {
			continue // the last element might be empty string
		}

		//nodeID
		nodeID, err := strconv.ParseInt(nodeElement, 10, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("decode nodeID failed from %s", nodeElement)
		}
		nodeIDs = append(nodeIDs, nodeID)
	}
	if len(nodeIDs) < 2 {
		return 0, nil, fmt.Errorf("too less nodeIDs %v from %v", nodeIDs, nodeElements)
	}

	return wayID, nodeIDs, nil
}

// ParseLine parses way2nodes csv format line.
func ParseLine(line string) (wayID int64, nodeIDs []int64, err error) {

	record := strings.Split(line, ",")
	return ParseRecord(record)
}

// FormatToRecord formats wayID and nodeIDs to csv record.
func FormatToRecord(wayID int64, nodeIDs []int64) []string {

	record := []string{}
	record = append(record, strconv.FormatUint(uint64(wayID), 10))
	for _, n := range nodeIDs {
		record = append(record, strconv.FormatUint(uint64(n), 10))
	}
	return record
}

// FormatToString formats wayID and nodeIDs to csv string line.
func FormatToString(wayID int64, nodeIDs []int64) string {

	record := FormatToRecord(wayID, nodeIDs)
	return strings.Join(record, ",")
}
