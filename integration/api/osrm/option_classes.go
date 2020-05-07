package osrm

import (
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
)

// OptionClasses represents OSRM exclude classes.
// https://github.com/Telenav/osrm-backend/blob/master/docs/http.md#general-options
type OptionClasses []string

// ParseOptionClasses parses OSRM option elements.
func ParseOptionClasses(s string) (OptionClasses, error) {

	s = strings.TrimSuffix(s, api.Comma) // remove the last `,` if exist

	classes := OptionClasses{}
	splits := strings.Split(s, api.Comma)
	for _, split := range splits {
		if len(split) == 0 {
			continue
		}
		classes = append(classes, split)
	}
	return classes, nil
}

func (o *OptionClasses) String() string {
	var s string
	for _, class := range *o {
		if len(class) == 0 {
			continue
		}

		if len(s) > 0 {
			s += api.Comma
		}
		s += class
	}
	return s
}
