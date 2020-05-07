package osrm

import (
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
)

// OptionElements represents values split by `;`, i.e. bearings, radiuses, hints, approaches
type OptionElements []string

// ParseOptionElemenets parses OSRM option elements.
func ParseOptionElemenets(s string) (OptionElements, error) {

	s = strings.TrimSuffix(s, api.Semicolon) // remove the last `;` if exist
	if len(s) == 0 {
		return OptionElements{}, nil
	}

	elements := OptionElements{}
	splits := strings.Split(s, api.Semicolon)
	for _, split := range splits {
		elements = append(elements, split)
	}
	return elements, nil
}

func (o *OptionElements) String() string {
	var s string
	for _, element := range *o {
		if len(s) > 0 {
			s += api.Semicolon
		}
		s += element
	}
	return s
}
