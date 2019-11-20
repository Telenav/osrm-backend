package osrmv1

// OSRM response code
// doc: https://github.com/Telenav/osrm-backend/blob/master-telenav/docs/http.md#code
const (
	CodeOK             = "Ok"             // Request could be processed as expected.
	CodeInvalidURL     = "InvalidUrl"     // URL string is invalid.
	CodeInvalidService = "InvalidService" // Service name is invalid.
	CodeInvalidVersion = "InvalidVersion" // Version is not found.
	CodeInvalidOptions = "InvalidOptions" // Options are invalid.
	CodeInvalidQuery   = "InvalidQuery"   // The query string is synctactically malformed.
	CodeInvalidValue   = "InvalidValue"   // The successfully parsed query parameters are invalid.
	CodeNoSegment      = "NoSegment"      // One of the supplied input coordinates could not snap to street segment.
	CodeTooBig         = "TooBig"         // The request size violates one of the service specific request size restrictions.
	CodeNoRouteFound   = "NoRoute"        // No route found.
)
