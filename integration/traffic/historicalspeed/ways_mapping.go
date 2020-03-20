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

// WaysMapping represents way->patterns(weekly),timezone,daylight_saving mapping for historical speeds querying.
type WaysMapping map[int64]*mappingItem // indexed by wayID: positive means forward, negative means backward

const (
	daysPerWeek = 7

	fieldsPerCSVLine                           = 9                    // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
	fieldsWithTimezoneDaylightSavingPerCSVLine = fieldsPerCSVLine + 2 // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S,TIME_ZONE,DAYLIGHT_SAVING
)

// Count returns how many ways(directed) mapping records.
func (w WaysMapping) Count() int {
	return len(w)
}

// Load loads data from csv files.
func (w *WaysMapping) Load(filesPath []string) error {

	for _, f := range filesPath {
		err := w.loadFromSingleFile(f)
		if err != nil {
			return err
		}
	}

	glog.Infof("Loaded way2patterns mapping count %d", w.Count())
	return nil
}

func (w *WaysMapping) loadFromSingleFile(filePath string) error {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", filePath)

	r := csv.NewReader(f)

	beforeLoadMappingCount := w.Count()
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

		wayID, mapping, err := parseWaysMappingRecord(record)
		if err != nil {
			if count == 0 {
				glog.V(2).Infof("Ignore head record %v due to parse failure: %v", record, err)
			} else {
				glog.Warningf("Parse record %v failed, err: %v", record, err)
			}
			continue
		}

		(*w)[wayID] = mapping
		count++
	}

	glog.V(1).Infof("Loaded way2patterns mapping from file %s, count %d, total succeed parsed count %d", filePath, w.Count()-beforeLoadMappingCount, count)
	return nil
}

func parseWaysMappingRecord(record []string) (int64, *mappingItem, error) {
	if len(record) != fieldsPerCSVLine && len(record) != fieldsWithTimezoneDaylightSavingPerCSVLine {
		return 0, nil, fmt.Errorf("expect %d or %d fields in csv record but got %d", fieldsPerCSVLine, fieldsWithTimezoneDaylightSavingPerCSVLine, len(record))
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

	if len(record) == fieldsWithTimezoneDaylightSavingPerCSVLine { // prase timezoen and daylight_saving if exist
		//TODO:
	}

	return wayID, &mapping, nil
}
