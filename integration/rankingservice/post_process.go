package rankingservice

import (
	"strconv"

	"github.com/golang/glog"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
)

func pickupRoutes(routes []*osrmv1.Route, num int) []*osrmv1.Route {
	if len(routes) <= num {
		return routes
	}
	return routes[:num]
}

func cleanupAnnotations(routes []*osrmv1.Route, annotations string) {
	if annotations != osrmv1.AnnotationsValueFalse {
		return // return all annotations even if want some
	}

	for _, route := range routes {
		for _, leg := range route.Legs {
			leg.Annotation = nil
		}
	}
}

func parseAlternativesNumber(alternatives string) int {
	if alternatives == osrmv1.AlternativesValueFalse {
		return 1
	} else if alternatives == osrmv1.AlternativesValueTrue {
		return 2
	}

	num, err := strconv.ParseUint(alternatives, 10, 32)
	if err != nil {
		glog.Warningf("parse alternatives %s failed, err %v", alternatives, err)
		return 1
	}
	return int(num)
}
