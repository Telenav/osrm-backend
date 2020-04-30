package options

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api"
	"github.com/golang/glog"
)

// ParseAlternatives parses route service Alternatives option.
func ParseAlternatives(s string) (string, int, error) {

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

// ParseSteps parses route service Steps option.
func ParseSteps(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		glog.Warning(err)
		return false, err
	}

	return b, nil
}

// ParseAnnotations parses route service Annotations option.
func ParseAnnotations(s string) (string, error) {

	validAnnotationsValues := map[string]struct{}{
		AnnotationsValueTrue:        struct{}{},
		AnnotationsValueFalse:       struct{}{},
		AnnotationsValueNodes:       struct{}{},
		AnnotationsValueDistance:    struct{}{},
		AnnotationsValueDuration:    struct{}{},
		AnnotationsValueDataSources: struct{}{},
		AnnotationsValueWeight:      struct{}{},
		AnnotationsValueSpeed:       struct{}{},
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

// ParseGeometries parses route service Geometries option.
func ParseGeometries(s string) (string, error) {
	validGeometriesValues := map[string]struct{}{
		GeometriesValuePolyline:  struct{}{},
		GeometriesValuePolyline6: struct{}{},
		GeometriesValueGeojson:   struct{}{},
	}

	if _, found := validGeometriesValues[s]; !found {
		err := fmt.Errorf("invalid geometries value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}

// ParseOverview parses route service Overview option.
func ParseOverview(s string) (string, error) {
	validOverviewValues := map[string]struct{}{
		OverviewValueSimplified: struct{}{},
		OverviewValueFull:       struct{}{},
		OverviewValueFalse:      struct{}{},
	}

	if _, found := validOverviewValues[s]; !found {
		err := fmt.Errorf("invalid overview value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}

// ParseContinueStraight parses route service ContinueStraight option.
func ParseContinueStraight(s string) (string, error) {
	validContinueStraightValues := map[string]struct{}{
		ContinueStraightValueDefault: struct{}{},
		ContinueStraightValueFalse:   struct{}{},
		ContinueStraightValueTrue:    struct{}{},
	}

	if _, found := validContinueStraightValues[s]; !found {
		err := fmt.Errorf("invalid continue_straight value: %s", s)
		glog.Warning(err)
		return "", err
	}
	return s, nil
}
