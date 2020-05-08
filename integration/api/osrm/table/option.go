package table

import (
	"fmt"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

// Table service Query Parameter/Option Keys
// https://github.com/Telenav/osrm-backend/blob/master/docs/http.md#table-service
const (
	OptionKeySources      = "sources"
	OptionKeyDestinations = "destinations"
	OptionKeyAnnotations  = "annotations"
)

// Annotations default value
const (
	OptionAnnotationsDefaultValue = OptionAnnotationsValueDuration
)

// Annotations values
const (
	OptionAnnotationsValueDistance = "distance"
	OptionAnnotationsValueDuration = "duration"
)

// ParseOptionAnnotations parses table service Annotations option.
func parseOptionAnnotations(s string) (string, error) {

	validAnnotationsValues := map[string]struct{}{
		OptionAnnotationsValueDistance: struct{}{},
		OptionAnnotationsValueDuration: struct{}{},
	}

	splits := strings.Split(s, api.Comma)
	for _, split := range splits {
		if _, found := validAnnotationsValues[split]; !found {

			err := fmt.Errorf("invalid annotations value: %s", s)
			glog.Warning(err)
			return "", err
		}
	}

	return s, nil
}
