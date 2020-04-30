// Package code defines response code of OSRM services.
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#code
package code

// OSRM Response codes
const (
	// Generic code
	// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#code
	OK             = "Ok"             // Request could be processed as expected.
	InvalidURL     = "InvalidUrl"     // URL string is invalid.
	InvalidService = "InvalidService" // Service name is invalid.
	InvalidVersion = "InvalidVersion" // Version is not found.
	InvalidOptions = "InvalidOptions" // Options are invalid.
	InvalidQuery   = "InvalidQuery"   // The query string is synctactically malformed.
	InvalidValue   = "InvalidValue"   // The successfully parsed query parameters are invalid.
	NoSegment      = "NoSegment"      // One of the supplied input coordinates could not snap to street segment.
	TooBig         = "TooBig"         // The request size violates one of the service specific request size restrictions.

	// Route service extra response code
	// https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#route-service
	NoRouteFound = "NoRoute" // No route found.
)
