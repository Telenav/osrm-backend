package rankingservice

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm"
)

func pickupRoutes(routes []*osrm.Route, num int) []*osrm.Route {
	if len(routes) <= num {
		return routes
	}
	return routes[:num]
}

func cleanupAnnotations(routes []*osrm.Route, annotations string) {
	if annotations != osrm.AnnotationsValueFalse {
		return // return all annotations even if want some
	}

	for _, route := range routes {
		for _, leg := range route.Legs {
			leg.Annotation = nil
		}
	}
}
