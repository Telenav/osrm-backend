package timezone

import (
	"fmt"
	"strconv"
)

// ParseTimezone parses timezone string to integer representation.
func ParseTimezone(s string) (int16, error) {

	timezone, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("parse timezone from %s failed, err %v", s, err)
	}
	if !isValidTimezone(int16(timezone)) {
		return 0, fmt.Errorf("parsed timezone %d from %s is invalid", timezone, s)
	}

	return int16(timezone), nil
}

// FormatTimezone converts timezone to string representation.
func FormatTimezone(timezone int16) string {

	if timezone >= 0 {
		return fmt.Sprintf("%03d", timezone) //e.g. 000
	}

	return fmt.Sprintf("%04d", timezone) //e.g. -070
}

const (
	minTimezone = -120
	maxTimezone = 140
)

func isValidTimezone(timezone int16) bool {
	if timezone < minTimezone || timezone > maxTimezone {
		return false
	}
	return true
}
