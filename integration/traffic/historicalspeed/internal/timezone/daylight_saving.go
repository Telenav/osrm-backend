package timezone

import (
	"fmt"
	"strconv"
)

// ParseDaylightSaving parses daylight saving string to integer representation.
func ParseDaylightSaving(s string) (int8, error) {

	dst, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("parse daylight saving from %s failed, err %v", s, err)
	}
	if !isValidDaylightSaving(int8(dst)) {
		return 0, fmt.Errorf("parsed daylight saving %d from %s is invalid", dst, s)
	}

	return int8(dst), nil
}

// FormatDaylightSaving converts daylight saving to string representation.
func FormatDaylightSaving(daylightSaving int8) string {

	return strconv.FormatInt(int64(daylightSaving), 10)
}

const (
	minDaylightSaving = 0
	maxDaylightSaving = 67
)

func isValidDaylightSaving(daylightSaving int8) bool {
	if daylightSaving < minDaylightSaving || daylightSaving > maxDaylightSaving {
		return false
	}
	return true
}
