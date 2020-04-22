package ranking

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/golang/glog"
)

func (h *Handler) retrieveWayIDs(routes []*route.Route) error {

	startTime := time.Now()

	var waysCount, nodesCount int
	for _, r := range routes {
		for _, l := range r.Legs {
			if l != nil && l.Annotation != nil {
				wayIDs, err := h.nodes2WayQuerier.QueryWays(l.Annotation.Nodes)
				if err != nil {
					return err
				}
				l.Annotation.Ways = wayIDs

				nodesCount += len(l.Annotation.Nodes)
				waysCount += len(l.Annotation.Ways)
			}
		}
	}

	glog.V(2).Infof("Retrieved %d wayIDs from %d nodeIDs, takes %f seconds.", waysCount, nodesCount, time.Now().Sub(startTime).Seconds())
	return nil
}
