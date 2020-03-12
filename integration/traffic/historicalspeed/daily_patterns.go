package historicalspeed

type dailyPattern [patternsPerDay]uint8

const (
	dailyPatternIntervalInMinutes = 15                                      // 15 minutes per value, e.g. speed_0 represents [00:00~00:15)
	patternsPerDay                = 24 * 60 / dailyPatternIntervalInMinutes // 96 per day
)
