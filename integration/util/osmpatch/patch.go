// Package osmpatch implements some utility to optimize handling for OSM PBF.
package osmpatch

import "github.com/qedus/osmpbf"

// IsValidWay returns whether the way from OSM is valid or not.
// This validation is for OSM PBF only.
// Rule: it is a valid/navigable way only if it has `highway` or `route` tag.
// https://github.com/Telenav/osrm-backend/blob/master/profiles/car.lua#L379
func IsValidWay(way *osmpbf.Way) bool {
	if way == nil {
		return false
	}

	_, hasHighwayTag := way.Tags["highway"]
	_, hasRouteTag := way.Tags["route"]

	if hasHighwayTag || hasRouteTag {
		return true
	}

	return false
}
