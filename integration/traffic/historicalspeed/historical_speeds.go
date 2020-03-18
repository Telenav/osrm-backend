package historicalspeed

import "time"

// Speeds stores historical speeds.
type Speeds struct {
	dailyPatterns       map[uint32]dailyPattern
	way2PatternsMapping map[int64]*mappingItem // indexed by wayID: positive means forward, negative means backward

	// allow multiple files
	dailyPatternsFilePath        []string
	ways2PatternsMappingFilePath []string
}

// New create a empty Speeds object.
func New(dailyPatternsFilePath, ways2PatternsMappingFilePath []string) *Speeds {

	return &Speeds{
		dailyPatterns:       map[uint32]dailyPattern{},
		way2PatternsMapping: map[int64]*mappingItem{},

		dailyPatternsFilePath:        dailyPatternsFilePath,
		ways2PatternsMappingFilePath: ways2PatternsMappingFilePath,
	}
}

// Load loads contents from file into memory.
func (s *Speeds) Load() error {
	if err := s.loadDailyPatterns(); err != nil {
		return err
	}
	if err := s.loadWaysPatternsMapping(); err != nil {
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
	return len(s.dailyPatterns)
}

// WaysCount returns how many directed ways have historical speeds.
func (s *Speeds) WaysCount() int {
	return len(s.way2PatternsMapping)
}
