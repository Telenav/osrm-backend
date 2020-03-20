package historicalspeed

import "time"

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
// wayID: positive means forward, negative means backward
// t: UTC time
func (s *Speeds) QueryHistoricalSpeed(wayID int64, t time.Time) float64 {
	//TODO:
	return 0.0
}

// DailyPatternsCount returns how many daily patterns.
func (s *Speeds) DailyPatternsCount() int {
	return s.dailyPatterns.count()
}

// WaysCount returns how many directed ways have historical speeds.
func (s *Speeds) WaysCount() int {
	return s.WaysMapping.Count()
}
