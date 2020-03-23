package historicalspeed

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Telenav/osrm-backend/integration/traffic/historicalspeed/internal/timezone"

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

	startTime := time.Now()

	for _, f := range filesPath {
		err := w.loadFromSingleFile(f)
		if err != nil {
			return err
		}
	}

	glog.Infof("Loaded way2patterns mapping count %d, takes %f seconds", w.Count(), time.Now().Sub(startTime).Seconds())
	return nil
}

// Dump dumps contents to csv
func (w *WaysMapping) Dump(filePath string, withCSVHeader bool) error {

	startTime := time.Now()

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
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

	glog.Infof("Dumped ways mapping to %s, csv header: %t, total count: %d, takes %f seconds", filePath, withCSVHeader, count, time.Now().Sub(startTime).Seconds())
	return nil
}

// UpdateTimezoneDaylightSaving updates timezone and daylight saving on way.
func (w *WaysMapping) UpdateTimezoneDaylightSaving(wayID int64, tzStr string, dstStr string) error {

	forwardItem, forwardOK := (*w)[wayID]
	backwardItem, backwardOK := (*w)[-wayID]
	if !forwardOK && !backwardOK {
		return fmt.Errorf("no historical speed on the way")
	}

	tz, err := timezone.ParseTimezone(tzStr)
	if err != nil {
		return err
	}
	dst, err := timezone.ParseDaylightSaving(dstStr)
	if err != nil {
		return err
	}

	if forwardOK {
		forwardItem.timezone = tz
		forwardItem.daylightSaving = dst
	}
	if backwardOK {
		backwardItem.timezone = tz
		backwardItem.daylightSaving = dst
	}

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
	r.ReuseRecord = true

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

		tz, err := timezone.ParseTimezone(record[9])
		if err != nil {
			return 0, nil, err
		}
		mapping.timezone = tz

		dst, err := timezone.ParseDaylightSaving(record[10])
		if err != nil {
			return 0, nil, err
		}
		mapping.daylightSaving = dst
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

	record = append(record, timezone.FormatTimezone(item.timezone))
	record = append(record, timezone.FormatDaylightSaving(item.daylightSaving))

	if len(record) != fieldsWithTimezoneDaylightSavingPerCSVLine {
		glog.Fatalf("expect record count %d but got %d", fieldsWithTimezoneDaylightSavingPerCSVLine, len(record))
	}
	return record
}
