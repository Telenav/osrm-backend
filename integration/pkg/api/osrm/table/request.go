package table

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"

// Request for OSRM table service
// http://project-osrm.org/docs/v5.5.1/api/#table-service
type Request struct {
	Service     string
	Version     string
	Profile     string
	Coordinates coordinate.Coordinates

	// options
	Sources      []*int
	Destinations []*int
	Annotations  string
}

// NewRequest create an empty table Request.
func NewRequest() *Request {

	return &Request{
		// Path
		Service:     "route",
		Version:     "v1",
		Profile:     "driving",
		Coordinates: coordinate.Coordinates{},
	}

}
