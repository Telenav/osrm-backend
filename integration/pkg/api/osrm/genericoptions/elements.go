package genericoptions

import (
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
)

// Elements represents values split by `;`, i.e. bearings, radiuses, hints, approaches
type Elements []string

// ParseElemenets parses OSRM option elements.
func ParseElemenets(s string) (Elements, error) {

	s = strings.TrimSuffix(s, api.Semicolon) // remove the last `;` if exist
	if len(s) == 0 {
		return Elements{}, nil
	}

	elements := Elements{}
	splits := strings.Split(s, api.Semicolon)
	for _, split := range splits {
		elements = append(elements, split)
	}
	return elements, nil
}

func (e *Elements) String() string {
	var s string
	for _, element := range *e {
		if len(s) > 0 {
			s += api.Semicolon
		}
		s += element
	}
	return s
}
