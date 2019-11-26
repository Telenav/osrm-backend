package rankingservice

import (
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
