// Package idsflag provides flag util for multiple IDs based on signed integer.
package idsflag

import (
	"strconv"
	"strings"
)

// IDs defines an IDs structure which also possible to work with flag package.
type IDs []int64

func (ids IDs) String() string {
	s := []string{}
	for _, id := range ids {
		s = append(s, strconv.FormatInt(id, 10))
	}
	return strings.Join(s, ",")
}

// Set sets the value of the named command-line flag.
func (ids *IDs) Set(value string) error {
	if len(value) == 0 {
		return nil
	}

	for _, v := range strings.Split(value, ",") {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			*ids = append(*ids, id)
		} else {
			return err
		}
	}
	return nil
}
