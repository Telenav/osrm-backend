package wayidsflag

import (
	"strconv"
	"strings"
)

// WayIDs defines an wayIDs structure which also possible to work with flag package.
type WayIDs []int64

func (w WayIDs) String() string {
	s := []string{}
	for _, wayID := range w {
		s = append(s, strconv.FormatInt(wayID, 10))
	}
	return strings.Join(s, ",")
}

// Set sets the value of the named command-line flag.
func (w *WayIDs) Set(value string) error {
	if len(value) == 0 {
		return nil
	}

	for _, way := range strings.Split(value, ",") {
		if wayID, err := strconv.ParseInt(way, 10, 64); err == nil {
			*w = append(*w, wayID)
		} else {
			return err
		}
	}
	return nil
}
