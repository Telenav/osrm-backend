package ranking

import (
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/golang/glog"
)

func (h *Handler) retrieveWayIDs(routes []*route.Route) error {

	startTime := time.Now()

	var legsCount, waysCount, nodesCount int
	for _, r := range routes {
		for _, l := range r.Legs {
			if l != nil && l.Annotation != nil {
				wayIDs, err := h.nodes2WayQuerier.QueryWays(l.Annotation.Nodes)
				if err != nil {
					return err
				}
				l.Annotation.Ways = wayIDs

				legsCount++
				nodesCount += len(l.Annotation.Nodes)
				waysCount += len(l.Annotation.Ways)
			}
		}
	}

	glog.V(2).Infof("Retrieved %d wayIDs from %d nodeIDs for %d legs(%d routes), takes %f seconds.", waysCount, nodesCount, legsCount, len(routes), time.Now().Sub(startTime).Seconds())
	return nil
}
