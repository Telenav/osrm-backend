package options

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

// Geometries values
const (
	GeometriesValuePolyline  = "polyline"
	GeometriesValuePolyline6 = "polyline6"
	GeometriesValueGeojson   = "geojson"

	GeometriesDefaultValue = GeometriesValuePolyline
)

// Overview values
const (
	OverviewValueSimplified = "simplified"
	OverviewValueFull       = "full"
	OverviewValueFalse      = "false"

	OverviewDefaultValue = OverviewValueSimplified
)

// ContinueStraight values
const (
	ContinueStraightValueDefault = "default"
	ContinueStraightValueTrue    = ValueTrue
	ContinueStraightValueFalse   = ValueFalse

	ContinueStraightDefaultValue = ContinueStraightValueDefault
)
