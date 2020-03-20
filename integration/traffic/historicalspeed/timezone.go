package historicalspeed

const (
	minTimezone = -120
	maxTimezone = 140

	minDaylightSaving = 0
	maxDaylightSaving = 67
)

func isValidTimezone(timezone int16) bool {
	if timezone < minTimezone || timezone > maxTimezone {
		return false
	}
	return true
}

func isValidDaylightSaving(daylightSaving int8) bool {
	if daylightSaving < minDaylightSaving || daylightSaving > maxDaylightSaving {
		return false
	}
	return true
}
