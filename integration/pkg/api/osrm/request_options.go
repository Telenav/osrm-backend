package osrm

// Generic Query Parameter/Option Keys
// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#requests
const (
	KeyBearings      = "bearings"       // {bearing};{bearing}[;{bearing} ...]
	KeyRadiuses      = "radiuses"       // {radius};{radius}[;{radius} ...]
	KeyGenerateHints = "generate_hints" // true(default), false
	KeyHints         = "hints"          // {hint};{hint}[;{hint} ...]
	KeyApproaches    = "approaches"     // {approach};{approach}[;{approach} ...]
	KeyExclude       = "exclude"        // {class}[,{class}]
)

// Common use choice values
const (
	ValueTrue  = "true"
	ValueFalse = "false"
)
