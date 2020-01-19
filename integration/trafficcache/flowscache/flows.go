package flowscache

import (
	"sync"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// Cache stores flows in memory.
type Cache struct {
	m     sync.RWMutex
	flows map[int64]*trafficproxy.Flow
}

// New creates a new cache to store flows in memory.
func New() *Cache {
	return &Cache{sync.RWMutex{}, map[int64]*trafficproxy.Flow{}}
}

//Clear clear the cache.
func (c *Cache) Clear() {
	c.m.Lock()
	defer c.m.Unlock()

	c.flows = map[int64]*trafficproxy.Flow{}
}

// Query returns Live Traffic Flow if exist.
func (c *Cache) Query(wayID int64) *trafficproxy.Flow {
	c.m.RLock()
	defer c.m.RUnlock()

	v, ok := c.flows[wayID]
	if ok {
		return v
	}
	return nil
}

// Count returns how many flows in the cache.
func (c *Cache) Count() int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return int64(len(c.flows))
}

// Update updates flows in cache.
func (c *Cache) Update(flowResp []*trafficproxy.FlowResponse) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, f := range flowResp {
		if f.Action == trafficproxy.Action_UPDATE {
			if inCacheFlow, ok := c.flows[f.Flow.WayID]; ok {
				if inCacheFlow.Timestamp <= f.Flow.Timestamp {
					c.flows[f.Flow.WayID] = f.Flow // use newer if exist
				}
				continue
			}
			c.flows[f.Flow.WayID] = f.Flow // store if not exist
			continue
		} else if f.Action == trafficproxy.Action_DELETE {
			delete(c.flows, f.Flow.WayID)
			continue
		}

		//undefined
		glog.Errorf("undefined flow action %d, flow %v", f.Action, f.Flow)
	}
}
