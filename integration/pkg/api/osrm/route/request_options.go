package route

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm"

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

// Alternatives values
const (
	AlternativesValueTrue  = osrm.ValueTrue
	AlternativesValueFalse = osrm.ValueFalse

	AlternativesDefaultValue = AlternativesValueFalse // default
)

// Steps values
const (
	StepsDefaultValue = false // default
)

// Annotations values
const (
	AnnotationsValueTrue        = osrm.ValueTrue
	AnnotationsValueFalse       = osrm.ValueFalse
	AnnotationsValueNodes       = "nodes"
	AnnotationsValueDistance    = "distance"
	AnnotationsValueDuration    = "duration"
	AnnotationsValueDataSources = "datasources"
	AnnotationsValueWeight      = "weight"
	AnnotationsValueSpeed       = "speed"

	AnnotationsDefaultValue = AnnotationsValueFalse // default
)
