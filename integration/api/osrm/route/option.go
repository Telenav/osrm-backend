package route

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

// ParseOptionAlternatives parses route service Alternatives option.
func parseOptionAlternatives(s string) (string, int, error) {

	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return s, int(n), nil
	}
	if b, err := strconv.ParseBool(s); err == nil {
		if b {
			return s, 2, nil // true : 2
		}
		return s, 1, nil // false : 1
	}

	err := fmt.Errorf("invalid alternatives value: %s", s)
	glog.Warning(err)
	return "", 1, err // use value 1 if fail
}

// ParseOptionSteps parses route service Steps option.
func parseOptionSteps(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		glog.Warning(err)
		return false, err
	}

	return b, nil
}

// ParseOptionAnnotations parses route service Annotations option.
func parseOptionAnnotations(s string) (string, error) {

	validAnnotationsValues := map[string]struct{}{
		OptionAnnotationsValueTrue:        struct{}{},
		OptionAnnotationsValueFalse:       struct{}{},
		OptionAnnotationsValueNodes:       struct{}{},
		OptionAnnotationsValueDistance:    struct{}{},
		OptionAnnotationsValueDuration:    struct{}{},
		OptionAnnotationsValueDataSources: struct{}{},
		OptionAnnotationsValueWeight:      struct{}{},
		OptionAnnotationsValueSpeed:       struct{}{},
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

// ParseOptionGeometries parses route service Geometries option.
func parseOptionGeometries(s string) (string, error) {
	validGeometriesValues := map[string]struct{}{
		OptionGeometriesValuePolyline:  struct{}{},
		OptionGeometriesValuePolyline6: struct{}{},
		OptionGeometriesValueGeojson:   struct{}{},
	}

	if _, found := validGeometriesValues[s]; !found {
		err := fmt.Errorf("invalid geometries value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}

// ParseOptionOverview parses route service Overview option.
func parseOptionOverview(s string) (string, error) {
	validOverviewValues := map[string]struct{}{
		OptionOverviewValueSimplified: struct{}{},
		OptionOverviewValueFull:       struct{}{},
		OptionOverviewValueFalse:      struct{}{},
	}

	if _, found := validOverviewValues[s]; !found {
		err := fmt.Errorf("invalid overview value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}

// ParseOptionContinueStraight parses route service ContinueStraight option.
func parseOptionContinueStraight(s string) (string, error) {
	validContinueStraightValues := map[string]struct{}{
		OptionContinueStraightValueDefault: struct{}{},
		OptionContinueStraightValueFalse:   struct{}{},
		OptionContinueStraightValueTrue:    struct{}{},
	}

	if _, found := validContinueStraightValues[s]; !found {
		err := fmt.Errorf("invalid continue_straight value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}
