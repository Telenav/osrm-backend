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

type way2PatternsMapping struct {
	m         map[int64]*mappingItem // indexed by wayID: positive means forward, negative means backward
	filesPath []string               // allow multiple files
}

const (
	daysPerWeek = 7

	fieldsPerMapping = 9 //LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
)

func (w way2PatternsMapping) count() int {
	return len(w.m)
}

func (w *way2PatternsMapping) load() error {

	for _, f := range w.filesPath {
		err := w.loadFromSingleFile(f)
		if err != nil {
			return err
		}
	}

	glog.Infof("Loaded way2patterns mapping count %d", w.count())
	return nil
}

func (w *way2PatternsMapping) loadFromSingleFile(filePath string) error {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", filePath)

	r := csv.NewReader(f)

	beforeLoadMappingCount := w.count()
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

		w.m[wayID] = mapping
		count++
	}

	glog.V(1).Infof("Loaded way2patterns mapping from file %s, count %d, total succeed parsed count %d", filePath, w.count()-beforeLoadMappingCount, count)
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
