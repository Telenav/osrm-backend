package oasis

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/api/oasis"
	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
)

// currRange: current energy level, represent by distance(unit: meters)
// destRange: energy level required to destination
// currRange - routeDistance > destRange means no charge station is needed
func hasEnoughEnergy(currRange, destRange float64, routeResp *route.Response) (bool, float64, error) {
	if len(routeResp.Routes) == 0 {
		err := fmt.Errorf("route response contains no route result")
		return false, 0, err
	}

	if currRange < 0 || destRange < 0 {
		err := fmt.Errorf("incorrect parameter of range")
		return false, 0, err
	}

	minRouteDistance := math.MaxFloat64

	for _, route := range routeResp.Routes {
		minRouteDistance = math.Min(minRouteDistance, route.Distance)
	}

	if minRouteDistance < 0 {
		err := fmt.Errorf("incorrect route distance in response with negative value")
		return false, 0, err
	}

	remainRange := currRange - minRouteDistance
	if remainRange > destRange {
		return true, remainRange, nil
	}
	return false, remainRange, nil
}

func generateOASISResponse4NoChargeNeeded(w http.ResponseWriter, routeResp *route.Response, remainRange float64) {
	w.WriteHeader(http.StatusOK)

	solution := new(oasis.Solution)
	solution.Distance = routeResp.Routes[0].Distance
	solution.Duration = routeResp.Routes[0].Duration
	solution.Weight = routeResp.Routes[0].Weight
	solution.RemainingRage = remainRange
	solution.WeightName = routeResp.Routes[0].WeightName

	r := new(oasis.Response)
	r.Code = "200"
	r.Message = "Success."
	r.Solutions = append(r.Solutions, solution)

	json.NewEncoder(w).Encode(r)
}
