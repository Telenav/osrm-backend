// Package genericoptions defines OSRM generic request options keys and values.
// https://github.com/Telenav/osrm-backend/blob/feature/fail-on-error/docs/http.md#requests
package genericoptions

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
