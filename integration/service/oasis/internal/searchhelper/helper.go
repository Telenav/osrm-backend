package searchhelper

import (
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/api/search/searchcoordinate"
)

// GenerateSearchRequest accepts center point and limitations and generate nearbychargestation.Request
func GenerateSearchRequest(location searchcoordinate.Coordinate, limit int, radius float64) (*nearbychargestation.Request, error) {
	req := nearbychargestation.NewRequest()
	req.Location = location
	if limit > 0 {
		req.Limit = limit
	}

	if radius > 0 {
		req.Radius = radius
	}

	return req, nil
}
