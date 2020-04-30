package options

import (
	"fmt"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

// ParseAnnotations parses table service Annotations option.
func ParseAnnotations(s string) (string, error) {

	validAnnotationsValues := map[string]struct{}{
		AnnotationsValueDistance: struct{}{},
		AnnotationsValueDuration: struct{}{},
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
