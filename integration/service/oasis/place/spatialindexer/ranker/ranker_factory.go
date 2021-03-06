package ranker

import (
	"github.com/Telenav/osrm-backend/integration/service/oasis/place"
	"github.com/Telenav/osrm-backend/integration/util/osrmconnector"
)

const (
	// SimpleRanker implements Raner's interface based on great circle distance
	SimpleRanker = "SimpleRanker"
	// OSRMBasedRanker implements Raner's interface based on OSRM
	OSRMBasedRanker = "OSRMBasedRanker"
)

// CreateRanker creates implementations of interface Ranker
func CreateRanker(rankerType string, oc *osrmconnector.OSRMConnector) place.Ranker {
	switch rankerType {
	case SimpleRanker:
		return newSimpleRanker()
	case OSRMBasedRanker:
		return newOsrmRanker(oc)
	default:
		return newSimpleRanker()
	}
}
