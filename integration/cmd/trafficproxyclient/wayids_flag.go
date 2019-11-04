package main

import (
	"strconv"
	"strings"
)

type wayIDsFlag struct {
	wayIDs []int64
}

func (w *wayIDsFlag) String() string {
	s := []string{}
	for _, wayID := range w.wayIDs {
		s = append(s, strconv.FormatInt(wayID, 10))
	}
	return strings.Join(s, ",")
}

func (w *wayIDsFlag) Set(value string) error {
	if len(value) == 0 {
		return nil
	}

	for _, way := range strings.Split(value, ",") {
		if wayID, err := strconv.ParseInt(way, 10, 64); err == nil {
			w.wayIDs = append(w.wayIDs, wayID)
		} else {
			return err
		}
	}
	return nil
}
