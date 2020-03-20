package historicalspeed

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type mappingItem struct {
	patternIDs     [daysPerWeek]uint32 // U,M,T,W,R,F,S
	timezone       int16
	daylightSaving int8
}

// WaysMapping represents way->patterns(weekly),timezone,daylight_saving mapping for historical speeds querying.
type WaysMapping map[int64]*mappingItem // indexed by wayID: positive means forward, negative means backward

const (
	daysPerWeek = 7

	fieldsPerCSVLine                           = 9                    // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S
	fieldsWithTimezoneDaylightSavingPerCSVLine = fieldsPerCSVLine + 2 // LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S,TIME_ZONE,DAYLIGHT_SAVING

	outputCSVHeader = "LINK_PVID,TRAVEL_DIRECTION,U,M,T,W,R,F,S,TIME_ZONE,DAYLIGHT_SAVING"
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

// Dump dumps contents to csv
func (w *WaysMapping) Dump(filePath string, withCSVHeader bool) error {

	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	defer f.Sync()
	glog.V(1).Infof("open %s succeed.\n", filePath)

	writer := csv.NewWriter(f)
	defer writer.Flush()
	if withCSVHeader {
		if err := writer.Write(strings.Split(outputCSVHeader, ",")); err != nil {
			return err
		}
	}

	var count int
	for k, v := range *w {
		if err := writer.Write(toWaysMappingRecord(k, v)); err != nil {
			return err
		}
		count++
	}

	glog.Infof("Dumped ways mapping to %s, csv header: %t, total count: %d", filePath, withCSVHeader, count)
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

	patternIDsRecord := record[2 : 2+daysPerWeek]
	mapping := mappingItem{}

	for i, v := range patternIDsRecord {
		patternID, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, nil, fmt.Errorf("parse patternID from %s failed, err %v", v, err)
		}
		mapping.patternIDs[i] = uint32(patternID)
	}

	if len(record) == fieldsWithTimezoneDaylightSavingPerCSVLine { // prase timezoen and daylight_saving if exist

		timezoneStr := record[9]
		timezone, err := strconv.ParseInt(timezoneStr, 10, 16)
		if err != nil {
			return 0, nil, fmt.Errorf("parse timezone from %s failed, err %v", timezoneStr, err)
		}
		if !isValidTimezone(int16(timezone)) {
			return 0, nil, fmt.Errorf("parsed timezone %d from %s is invalid", timezone, timezoneStr)
		}
		mapping.timezone = int16(timezone)

		dstStr := record[10]
		dst, err := strconv.ParseInt(dstStr, 10, 8)
		if err != nil {
			return 0, nil, fmt.Errorf("parse daylight saving from %s failed, err %v", dstStr, err)
		}
		if !isValidDaylightSaving(int8(dst)) {
			return 0, nil, fmt.Errorf("parsed daylight saving %d from %s is invalid", dst, dstStr)
		}
		mapping.daylightSaving = int8(dst)
	}

	return wayID, &mapping, nil
}

func toWaysMappingRecord(wayID int64, item *mappingItem) []string {
	record := []string{}
	if wayID >= 0 {
		record = append(record, strconv.FormatInt(wayID, 10), "T")
	} else {
		record = append(record, strconv.FormatInt(-wayID, 10), "F")
	}

	for _, v := range item.patternIDs {
		record = append(record, strconv.FormatUint(uint64(v), 10))
	}

	if item.timezone >= 0 {
		record = append(record, fmt.Sprintf("%03d", item.timezone)) //e.g. 000
	} else {
		record = append(record, fmt.Sprintf("%04d", item.timezone)) //e.g. -070
	}

	record = append(record, strconv.FormatInt(int64(item.daylightSaving), 10))

	if len(record) != fieldsWithTimezoneDaylightSavingPerCSVLine {
		glog.Fatalf("expect record count %d but got %d", fieldsWithTimezoneDaylightSavingPerCSVLine, len(record))
	}
	return record
}
