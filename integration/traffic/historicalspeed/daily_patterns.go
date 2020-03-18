package historicalspeed

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/golang/glog"
)

type dailyPattern []uint8

const (
	dailyPatternIntervalInMinutes = 15                                      // 15 minutes per value, e.g. speed_0 represents [00:00~00:15)
	patternsPerDay                = 24 * 60 / dailyPatternIntervalInMinutes // 96 per day

	fieldsPerPattern = 1 + patternsPerDay // patternID and patterns
)

func (s *Speeds) loadDailyPatterns() error {

	for _, f := range s.dailyPatternsFilePath {
		err := s.loadDailyPatternsFromSingleFile(f)
		if err != nil {
			return err
		}
	}

	glog.Infof("Loaded daily patterns count %d", len(s.dailyPatterns))
	return nil
}

func (s *Speeds) loadDailyPatternsFromSingleFile(filePath string) error {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", filePath)

	r := csv.NewReader(f)

	beforeLoadPatternsCount := len(s.dailyPatterns)
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

		id, pattern, err := parseDailyPatternRecord(record)
		if err != nil {
			if count == 0 {
				glog.V(2).Infof("Ignore head record %v due to parse failure: %v", record, err)
			} else {
				glog.Warningf("Parse record %v failed, err: %v", record, err)
			}
			continue
		}

		s.dailyPatterns[id] = pattern
		count++
	}

	glog.V(1).Infof("Loaded daily patterns from file %s, count %d, total succeed parsed count %d", filePath, len(s.dailyPatterns)-beforeLoadPatternsCount, count)
	return nil
}

func parseDailyPatternRecord(record []string) (uint32, dailyPattern, error) {

	if len(record) != fieldsPerPattern {
		return 0, dailyPattern{}, fmt.Errorf("expect %d fields in csv record but got %d", fieldsPerPattern, len(record))
	}

	patternID, err := strconv.ParseUint(record[0], 10, 32)
	if err != nil {
		return 0, dailyPattern{}, fmt.Errorf("parse daily patternID from %s error: %v", record[0], err)
	}

	speedsRecord := record[1:]
	pattern := make(dailyPattern, patternsPerDay, patternsPerDay)
	for i, v := range speedsRecord {
		speed, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, dailyPattern{}, fmt.Errorf("parse speed from %s failed, err %v", v, err)
		}
		if speed > math.MaxUint8 {
			return 0, dailyPattern{}, fmt.Errorf("parsed invalid speed %d (> %d) from %s", speed, math.MaxUint8, v)
		}
		pattern[i] = uint8(speed)
	}

	return uint32(patternID), pattern, nil
}
