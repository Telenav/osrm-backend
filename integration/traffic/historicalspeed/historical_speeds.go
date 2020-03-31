package historicalspeed

import (
	"time"

	"github.com/golang/glog"
)

// Speeds stores historical speeds.
type Speeds struct {
	dailyPatterns
	WaysMapping

	// allow multiple files
	dailyPatternsFilePath []string
	waysMappingFilePath   []string
}

// New create a empty Speeds object.
func New(dailyPatternsFilePath, waysMappingFilePath []string) *Speeds {

	return &Speeds{
		dailyPatterns: dailyPatterns{},
		WaysMapping:   WaysMapping{},

		dailyPatternsFilePath: dailyPatternsFilePath,
		waysMappingFilePath:   waysMappingFilePath,
	}
}

// Load loads contents from file into memory.
func (s *Speeds) Load() error {
	if err := s.dailyPatterns.load(s.dailyPatternsFilePath); err != nil {
		return err
	}
	if err := s.WaysMapping.Load(s.waysMappingFilePath); err != nil {
		return err
	}

	return nil
}

// QueryHistoricalSpeed return the speed for a way at a specified time
// wayID: positive means travel forward, which is travel following edge's point sequence, negative means backward
// t: UTC time
func (s *Speeds) QueryHistoricalSpeed(wayID int64, t time.Time) (float64, bool) {

	m, ok := s.WaysMapping[wayID]
	if !ok {
		return 0, false
	}

	//TODO: t from UTC to way's Local

	// find dailyPatternID
	patternID := m.getDailyPatternID(t.Weekday())

	// find pattern
	pattern, ok := s.dailyPatterns[patternID]
	if !ok {
		glog.Fatalf("Can not find daily pattern for PatternID %d on wayID %d.", patternID, wayID)
		return 0, false
	}

	return float64(pattern.querySpeed(t)), true
}

// DailyPatternsCount returns how many daily patterns.
func (s *Speeds) DailyPatternsCount() int {
	return s.dailyPatterns.count()
}

// WaysCount returns how many directed ways have historical speeds.
func (s *Speeds) WaysCount() int {
	return s.WaysMapping.Count()
}
