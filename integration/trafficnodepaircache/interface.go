package trafficnodepaircache

import (
	proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
)

// Cache is used to cache live traffic and provide query interfaces.
type Cache struct {
}

// New creates a new Cache instance.
func New() *Cache {
	c := Cache{}
	return &c
}

// Eat implements trafficproxyclient.Eater inteface.
func (c *Cache) Eat(r proxy.TrafficResponse) {
	//TODO:
}

// QueryFlow returns Live Traffic Flow if exist.
func (c *Cache) QueryFlow(fromNode, toNode int64) *proxy.Flow {

	//TODO:
	return nil
}

// QueryIncidents returns Live Traffic Incident if exist.
// Be aware that there might be more than one incident on one fromNode-toNode.
func (c *Cache) QueryIncidents(fromNode, toNode int64) []*proxy.Incident {

	//TODO:
	return nil
}

// IsBlockedByIncident checks whether the way has blocking incident.
func (c *Cache) IsBlockedByIncident(fromNode, toNode int64) bool {

	incidents := c.QueryIncidents(fromNode, toNode)
	for _, incident := range incidents {
		if incident != nil && incident.IsBlocking {
			return true
		}
	}
	return false
}
