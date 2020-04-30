// Package options defines OSRM route service request options.
package options

// Route service Query Parameter/Option Keys
// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#route-service
const (
	KeyAlternatives     = "alternatives"      // true, false(default) or number
	KeySteps            = "steps"             // true, false(default)
	KeyAnnotations      = "annotations"       // true, false(default), nodes, distance, duration, datasources, weight, speed
	KeyGeometries       = "geometries"        // polyline(default), polyline6, geojson
	KeyOverview         = "overview"          // simplified(default), full, false
	KeyContinueStraight = "continue_straight" // default(default), true, false
	KeyWaypoints        = "waypoints"         // {index};{index};{index}...

)
