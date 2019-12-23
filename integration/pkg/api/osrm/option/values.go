package option

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
