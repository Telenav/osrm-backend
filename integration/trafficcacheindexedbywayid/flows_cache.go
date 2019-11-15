package trafficcacheindexedbywayid

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

func (c *Cache) queryFlow(wayID int64) *proxy.Flow {
	v, ok := c.flowsCache.Load(wayID)
	if ok {
		return v.(*proxy.Flow)
	}
	return nil
}

func (c *Cache) updateFlows(flowResp []*proxy.FlowResponse) {
	for _, f := range flowResp {
		if f.Action == proxy.Action_UPDATE || f.Action == proxy.Action_ADD { //TODO: Action_ADD will be removed soon
			c.flowsCache.Store(f.Flow.WayId, f.Flow)
			continue
		} else if f.Action == proxy.Action_DELETE {
			c.flowsCache.Delete(f.Flow.WayId)
			continue
		}

		//undefined
		glog.Errorf("undefined flow action %d, flow %v", f.Action, f.Flow)
	}
}

func (c *Cache) clearFlows() {
	c.flowsCache.Range(func(key interface{}, value interface{}) bool {
		c.flowsCache.Delete(key)
		return true
	})
}

func (c *Cache) flowCount() int64 {
	var count int64
	c.flowsCache.Range(func(key interface{}, value interface{}) bool {
		count++
		return true
	})
	return count
}
