package searchconnector

import (
	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
)

// TNSearchConnector wraps the communication with TN search server
type TNSearchConnector struct {
	searchClient *tnSearchHTTPClient
}

// NewTNSearchConnector creates TNSearchConnector object
func NewTNSearchConnector(searchEndpoint, apiKey, apiSignature string) *TNSearchConnector {
	search := &TNSearchConnector{
		searchClient: newTNSearchHTTPClient(searchEndpoint, apiKey, apiSignature),
	}
	go search.searchClient.start()
	return search
}

// ChargeStationSearch returns a channel immediately.  Response information could be retrieved from the channel when ready.
func (sc *TNSearchConnector) ChargeStationSearch(req *nearbychargestation.Request) <-chan ChargeStationsResponse {
	return sc.searchClient.submitSearchReq(req)
}

// Stop will stop TNSearchConnector
func (sc *TNSearchConnector) Stop() {
	// todo
}
