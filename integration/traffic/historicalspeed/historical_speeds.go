package historicalspeed

import "time"

// Speeds stores historical speeds.
type Speeds struct {
	dailyPatterns       map[uint32]dailyPattern
	way2PatternsMapping map[int64]mappingItem // indexed by wayID: positive means forward, negative means backward

	dailyPatternsFilePath        string
	ways2PatternsMappingFilePath string
}

// New create a empty Speeds object.
func New(dailyPatternsFilePath, ways2PatternsMappingFilePath string) *Speeds {

	return &Speeds{
		dailyPatterns:       map[uint32]dailyPattern{},
		way2PatternsMapping: map[int64]mappingItem{},

		dailyPatternsFilePath:        dailyPatternsFilePath,
		ways2PatternsMappingFilePath: ways2PatternsMappingFilePath,
	}
}

// Load loads contents from file into memory.
func (s *Speeds) Load() error {
	// TODO:
	return nil
}

// QueryHistoricalSpeed return the speed for a way at a specified time
// wayID: positive means forward, negative means backward
// t: UTC time
func (s Speeds) QueryHistoricalSpeed(wayID int64, t time.Time) float64 {
	//TODO:
	return 0.0
}
