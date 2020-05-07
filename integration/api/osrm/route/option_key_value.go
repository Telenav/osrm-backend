package route

// Route service Query Parameter/Option Keys
// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#route-service
const (
	OptionKeyAlternatives     = "alternatives"      // true, false(default) or number
	OptionKeySteps            = "steps"             // true, false(default)
	OptionKeyAnnotations      = "annotations"       // true, false(default), nodes, distance, duration, datasources, weight, speed
	OptionKeyGeometries       = "geometries"        // polyline(default), polyline6, geojson
	OptionKeyOverview         = "overview"          // simplified(default), full, false
	OptionKeyContinueStraight = "continue_straight" // default(default), true, false
	OptionKeyWaypoints        = "waypoints"         // {index};{index};{index}...
)

// Common use choice values
const (
	OptionValueTrue  = "true"
	OptionValueFalse = "false"
)

// Alternatives values
const (
	OptionAlternativesValueTrue  = OptionValueTrue
	OptionAlternativesValueFalse = OptionValueFalse

	OptionAlternativesDefaultValue = OptionAlternativesValueFalse // default
)

// Steps values
const (
	OptionStepsDefaultValue = false // default
)

// Annotations values
const (
	OptionAnnotationsValueTrue        = OptionValueTrue
	OptionAnnotationsValueFalse       = OptionValueFalse
	OptionAnnotationsValueNodes       = "nodes"
	OptionAnnotationsValueDistance    = "distance"
	OptionAnnotationsValueDuration    = "duration"
	OptionAnnotationsValueDataSources = "datasources"
	OptionAnnotationsValueWeight      = "weight"
	OptionAnnotationsValueSpeed       = "speed"

	OptionAnnotationsDefaultValue = OptionAnnotationsValueFalse // default
)

// Geometries values
const (
	OptionGeometriesValuePolyline  = "polyline"
	OptionGeometriesValuePolyline6 = "polyline6"
	OptionGeometriesValueGeojson   = "geojson"

	OptionGeometriesDefaultValue = OptionGeometriesValuePolyline
)

// Overview values
const (
	OptionOverviewValueSimplified = "simplified"
	OptionOverviewValueFull       = "full"
	OptionOverviewValueFalse      = "false"

	OptionOverviewDefaultValue = OptionOverviewValueSimplified
)

// ContinueStraight values
const (
	OptionContinueStraightValueDefault = "default"
	OptionContinueStraightValueTrue    = OptionValueTrue
	OptionContinueStraightValueFalse   = OptionValueFalse

	OptionContinueStraightDefaultValue = OptionContinueStraightValueDefault
)
