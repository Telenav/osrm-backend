package ranking

import (
	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

func (h *Handler) retrieveWayIDs(routes []*route.Route) error {

	for _, r := range routes {
		for _, l := range r.Legs {
			if l != nil && l.Annotation != nil {
				wayIDs, err := h.nodes2WayQuerier.QueryWays(l.Annotation.Nodes)
				if err != nil {
					return err
				}
				l.Annotation.Ways = wayIDs
			}
		}
	}
	return nil
}
