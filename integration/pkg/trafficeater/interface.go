package trafficeater

import proxy "github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"

// Eater is the interface that wraps the basic Eat method.
type Eater interface {

	// Eat consumes traffic responses.
	Eat(proxy.TrafficResponse)
}
