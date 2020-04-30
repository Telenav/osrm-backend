package coordinate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

// Index indicates which NO. of Coordinates.
type Index uint

// Indexes represents a list of Index.
type Indexes []Index

func (i *Indexes) String() string {
	var s string
	for _, v := range *i {
		if len(s) > 0 {
			s += api.Semicolon
		}
		s += strconv.FormatUint(uint64(v), 10)
	}
	return s
}

// PraseIndexes parses string to indexes of coordinates.
func PraseIndexes(s string) (Indexes, error) {
	indexes := []Index{}

	splits := strings.Split(s, api.Semicolon)
	for _, split := range splits {
		if len(split) == 0 {
			continue
		}
		n, err := strconv.ParseUint(split, 10, 32)
		if err != nil {
			fullErr := fmt.Errorf("invalid indexes value: %s, err %v", s, err)
			glog.Warning(fullErr)
			return nil, fullErr
		}
		indexes = append(indexes, Index(n))
	}
	return indexes, nil
}
