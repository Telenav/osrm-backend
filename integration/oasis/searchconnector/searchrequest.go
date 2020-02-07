package searchconnector

import "github.com/Telenav/osrm-backend/integration/oasis"

type request struct {
	url         string
	searchRespC chan oasis.ChargeStationsResponse
}

func newTNSearchRequest(url string) *request {
	return &request{
		url:         url,
		searchRespC: make(chan oasis.ChargeStationsResponse),
	}
}
