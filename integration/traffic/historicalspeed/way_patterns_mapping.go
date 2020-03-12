package historicalspeed

type mappingItem struct {
	patternIDs     [daysPerWeek]uint32 // U,M,T,W,R,F,S
	timezone       int16
	daylightSaving int8
}

const (
	daysPerWeek = 7
)
