package genericoptions

import (
	"strconv"

	"github.com/golang/glog"
)

// ParseGenerateHints parses generic generate_hints option.
func ParseGenerateHints(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		glog.Warning(err)
		return false, err
	}

	return b, nil
}
