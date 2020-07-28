package searchconnector

type request struct {
	url         string
	searchRespC chan ChargeStationsResponse
}

func newTNSearchRequest(url string) *request {
	return &request{
		url:         url,
		searchRespC: make(chan ChargeStationsResponse),
	}
}
