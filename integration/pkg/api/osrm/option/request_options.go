// Package option defines OSRM request options keys and values.
// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#requests
package option

// OSRM request option keys
const (
	// Generic Query Parameter/Option Keys
	// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#requests
	KeyBearings      = "bearings"       // {bearing};{bearing}[;{bearing} ...]
	KeyRadiuses      = "radiuses"       // {radius};{radius}[;{radius} ...]
	KeyGenerateHints = "generate_hints" // true(default), false
	KeyHints         = "hints"          // {hint};{hint}[;{hint} ...]
	KeyApproaches    = "approaches"     // {approach};{approach}[;{approach} ...]
	KeyExclude       = "exclude"        // {class}[,{class}]

	// Route service Query Parameter/Option Keys
	// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#route-service
	KeyAlternatives     = "alternatives"      // true, false(default) or number
	KeySteps            = "steps"             // true, false(default)
	KeyAnnotations      = "annotations"       // true, false(default), nodes, distance, duration, datasources, weight, speed
	KeyGeometries       = "geometries"        // polyline(default), polyline6, geojson
	KeyOverview         = "overview"          // simplified(default), full, false
	KeyContinueStraight = "continue_straight" // default(default), true, false
	KeyWaypoints        = "waypoints"         // {index};{index};{index}...

)

// Common use choice values
const (
	ValueTrue  = "true"
	ValueFalse = "false"
)

// Alternatives values
const (
	AlternativesValueTrue  = ValueTrue
	AlternativesValueFalse = ValueFalse

	AlternativesDefaultValue = AlternativesValueFalse // default
)

// Steps values
const (
	StepsDefaultValue = false // default
)

// Annotations values
const (
	AnnotationsValueTrue        = ValueTrue
	AnnotationsValueFalse       = ValueFalse
	AnnotationsValueNodes       = "nodes"
	AnnotationsValueDistance    = "distance"
	AnnotationsValueDuration    = "duration"
	AnnotationsValueDataSources = "datasources"
	AnnotationsValueWeight      = "weight"
	AnnotationsValueSpeed       = "speed"

	AnnotationsDefaultValue = AnnotationsValueFalse // default
)
