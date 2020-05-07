package osrm

import (
	"strconv"

	"github.com/golang/glog"
)

// Generic Query Parameter/Option Keys
// https://github.com/Telenav/osrm-backend/blob/master/docs/http.md#general-options
const (
	OptionKeyBearings      = "bearings"       // {bearing};{bearing}[;{bearing} ...]
	OptionKeyRadiuses      = "radiuses"       // {radius};{radius}[;{radius} ...]
	OptionKeyGenerateHints = "generate_hints" // true(default), false
	OptionKeyHints         = "hints"          // {hint};{hint}[;{hint} ...]
	OptionKeyApproaches    = "approaches"     // {approach};{approach}[;{approach} ...]
	OptionKeyExclude       = "exclude"        // {class}[,{class}]
)

// GenerateHints values
const (
	OptionGenerateHintsDefaultValue = true // default
)

// ParseOptionGenerateHints parses generic generate_hints option.
func ParseOptionGenerateHints(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		glog.Warning(err)
		return false, err
	}

	return b, nil
}
