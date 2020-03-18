package historicalspeed

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/golang/glog"
)

type mappingItem struct {
	patternIDs     [daysPerWeek]uint32 // U,M,T,W,R,F,S
	timezone       int16               //TODO: timezone
	daylightSaving int8                //TODO: daylight saving
}

const (
	daysPerWeek = 7

	fieldsPerMapping = 9 //LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
)

func (s *Speeds) loadWaysPatternsMapping() error {

	for _, f := range s.ways2PatternsMappingFilePath {
		err := s.loadWaysPatternsMappingFromSingleFile(f)
		if err != nil {
			return err
		}
	}

	glog.Infof("Loaded way2patterns mapping count %d", len(s.way2PatternsMapping))
	return nil
}

func (s *Speeds) loadWaysPatternsMappingFromSingleFile(filePath string) error {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", filePath)

	r := csv.NewReader(f)

	beforeLoadMappingCount := len(s.way2PatternsMapping)
	var count int // succeed parsed count
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			return err
		}

		wayID, mapping, err := parseWay2PatternsMapping(record)
		if err != nil {
			if count == 0 {
				glog.V(2).Infof("Ignore head record %v due to parse failure: %v", record, err)
			} else {
				glog.Warningf("Parse record %v failed, err: %v", record, err)
			}
			continue
		}

		s.way2PatternsMapping[wayID] = mapping
		count++
	}

	glog.V(1).Infof("Loaded way2patterns mapping from file %s, count %d, total succeed parsed count %d", filePath, len(s.way2PatternsMapping)-beforeLoadMappingCount, count)
	return nil
}

func parseWay2PatternsMapping(record []string) (int64, *mappingItem, error) {
	if len(record) != fieldsPerMapping {
		return 0, nil, fmt.Errorf("expect %d fields in csv record but got %d", fieldsPerMapping, len(record))
	}

	undirectedWayID, err := strconv.ParseUint(record[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("parse wayID from %s error: %v", record[0], err)
	}
	wayID := int64(undirectedWayID)

	forward, err := strconv.ParseBool(record[1])
	if err != nil {
		return 0, nil, fmt.Errorf("parse wayID direction from %s error: %v", record[1], err)
	}
	if !forward {
		wayID = -wayID
	}

	patternIDsRecord := record[2:]
	mapping := mappingItem{}

	for i, v := range patternIDsRecord {
		patternID, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, nil, fmt.Errorf("parse patternID from %s failed, err %v", v, err)
		}
		mapping.patternIDs[i] = uint32(patternID)
	}

	return wayID, &mapping, nil
}
