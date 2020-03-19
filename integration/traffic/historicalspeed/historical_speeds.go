package historicalspeed

import "time"

// Speeds stores historical speeds.
type Speeds struct {
	dailyPatterns
	way2PatternsMapping
}

// New create a empty Speeds object.
func New(dailyPatternsFilePath, ways2PatternsMappingFilePath []string) *Speeds {

	return &Speeds{
		dailyPatterns{
			map[uint32]dailyPattern{}, dailyPatternsFilePath,
		},
		way2PatternsMapping{
			map[int64]*mappingItem{}, ways2PatternsMappingFilePath,
		},
	}
}

// Load loads contents from file into memory.
func (s *Speeds) Load() error {
	if err := s.dailyPatterns.load(); err != nil {
		return err
	}
	if err := s.way2PatternsMapping.load(); err != nil {
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
	return s.way2PatternsMapping.count()
}
