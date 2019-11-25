package rankingservice

import "github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"

// ranking routes and pick up the best alternatives to return.
func (h *Handler) ranking(routes []*osrmv1.Route, alternatives int) []*osrmv1.Route {
	if len(routes) == 0 || alternatives <= 0 {
		return routes
	}

	//TODO:

	return routes
}
