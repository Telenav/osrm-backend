package genericoptions

import (
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
)

// Classes represents OSRM exclude classes.
// https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#general-options
type Classes []string

// ParseClasses parses OSRM option elements.
func ParseClasses(s string) (Classes, error) {

	s = strings.TrimSuffix(s, api.Comma) // remove the last `,` if exist

	classes := Classes{}
	splits := strings.Split(s, api.Comma)
	for _, split := range splits {
		if len(split) == 0 {
			continue
		}
		classes = append(classes, split)
	}
	return classes, nil
}

func (c *Classes) String() string {
	var s string
	for _, class := range *c {
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
