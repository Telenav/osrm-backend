package osrmhelper

import (
	"fmt"
	"strconv"

	"github.com/Telenav/osrm-backend/integration/pkg/api"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/coordinate"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route/options"
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/table"
)

// GenerateTableReq4Points accept two group of points and generate osrm table request
func GenerateTableReq4Points(startPoints coordinate.Coordinates, endPoints coordinate.Coordinates) (*table.Request, error) {
	if len(startPoints) == 0 || len(endPoints) == 0 {
		return nil, fmt.Errorf("calling function with empty points")
	}

	// generate table request
	req := table.NewRequest()
	req.Coordinates = append(startPoints, endPoints...)

	count := 0
	for i := range startPoints {
		str := strconv.Itoa(i)
		req.Sources = append(req.Sources, str)
		count++
	}
	for i := range endPoints {
		str := strconv.Itoa(i + count)
		req.Destinations = append(req.Destinations, str)
	}

	req.Annotations = options.AnnotationsValueDistance + api.Comma + options.AnnotationsValueDuration
	return req, nil
}
