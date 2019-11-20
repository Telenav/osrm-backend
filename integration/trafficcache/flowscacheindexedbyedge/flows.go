package flowscacheindexedbyedge

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/wayidsmap"

	"github.com/Telenav/osrm-backend/integration/graph"
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// Cache stores flows in memory.
type Cache struct {
	m              sync.RWMutex
	flows          map[graph.Edge]*proxy.Flow
	affectedWayIDs map[int64]struct{}
	wayID2Edges    wayidsmap.Way2Edges
}

// New creates a new cache to store flows in memory.
func New(wayID2Edges wayidsmap.Way2Edges) *Cache {
	if wayID2Edges == nil {
		glog.Fatal("empty wayID2Edges")
		return nil
	}

	return &Cache{sync.RWMutex{},
		map[graph.Edge]*proxy.Flow{},
		map[int64]struct{}{},
		wayID2Edges}
}

//Clear clear the cache.
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.flows = map[graph.Edge]*proxy.Flow{}
	c.affectedWayIDs = map[int64]struct{}{}
}

// QueryEdge returns Live Traffic Flow for Edge if exist.
func (c *Cache) QueryEdge(edge graph.Edge) *proxy.Flow {
	c.m.RLock()
	defer c.m.RUnlock()

	v, ok := c.flows[edge]
	if ok {
		return v
	}
	return nil
}

// QueryEdges returns Live Traffic Flows for Edges if exist.
func (c *Cache) QueryEdges(edges []graph.Edge) []*proxy.Flow {
	c.m.RLock()
	defer c.m.RUnlock()

	out := make([]*proxy.Flow, len(edges), len(edges))
	for i := range edges {
		v, ok := c.flows[edges[i]]
		if ok {
			out[i] = v
			continue
		}
		out[i] = nil
	}
	return out
}

// Count returns how many flows in the cache.
func (c *Cache) Count() int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return int64(len(c.flows))
}

// AffectedWaysCount returns how many ways affected by these flows in the cache.
func (c *Cache) AffectedWaysCount() int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return int64(len(c.affectedWayIDs))
}

// Update updates flows in cache.
func (c *Cache) Update(flowResp []*proxy.FlowResponse) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, f := range flowResp {
		if f.Action == proxy.Action_UPDATE || f.Action == proxy.Action_ADD { //TODO: Action_ADD will be removed soon
			edges := c.wayID2Edges.WayID2Edges(f.Flow.WayId)
			for _, e := range edges {
				c.flows[e] = f.Flow
			}
			c.affectedWayIDs[f.Flow.WayId] = struct{}{}
			continue
		} else if f.Action == proxy.Action_DELETE {
			edges := c.wayID2Edges.WayID2Edges(f.Flow.WayId)
			for _, e := range edges {
				delete(c.flows, e)
			}
			delete(c.affectedWayIDs, f.Flow.WayId)
			continue
		}

		//undefined
		glog.Errorf("undefined flow action %d, flow %v", f.Action, f.Flow)
	}
}
