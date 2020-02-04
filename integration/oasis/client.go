package oasis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/backend"
)

// todo
// Merge RouteRequest/TableRequest to Request
// Handle initial check
// Add unit test

type RouteRequest struct {
	url      string
	respChan chan RouteResponse
}

type TableRequest struct {
	url string
}

func NewRouteRequest(url string) *RouteRequest {
	return &RouteRequest{
		url:      url,
		respChan: make(chan RouteResponse),
	}
}

type RouteResponse struct {
	Resp *route.Response
	Err  error
}

type OsrmClient struct {
	osrmBackendEndpoint string
	httpclient          http.Client
	routeRequest        chan RouteRequest
	tableRequest        chan TableRequest
}

func NewOsrmClient(osrmEndpoint string) *OsrmClient {
	osrmclient := &OsrmClient{
		osrmBackendEndpoint: osrmEndpoint,
		httpclient:          http.Client{Timeout: backend.Timeout()},
		routeRequest:        make(chan RouteRequest),
		tableRequest:        make(chan TableRequest),
	}

	// todo: make sure endpoint is valid
	go osrmclient.handle()

	return osrmclient
}

func (oc *OsrmClient) Request4Route() <-chan RouteResponse {
	//response := make(chan RouteResponse)
	url := "http://internal-a50649fcb01f011ea84ac0603197f379-25320380.us-west-2.elb.amazonaws.com:5001/route/v1/driving/-122.006349,37.364336;-121.875654,37.313767"
	request := *NewRouteRequest(url)
	oc.routeRequest <- request
	return request.respChan
}

type message struct {
	resp *http.Response
	err  error
	r    chan RouteResponse
}

func (oc *OsrmClient) handleRoute(routeReq RouteRequest, routeChan chan<- message) {
	fmt.Println("+++ " + routeReq.url)
	resp, err := oc.httpclient.Get(routeReq.url)

	m := message{
		resp: resp,
		err:  err,
		r:    routeReq.respChan,
	}
	routeChan <- m
}

func handleRouteResponse(m *message) {
	resp := m.resp
	defer resp.Body.Close()
	var response RouteResponse

	if m.err != nil {
		fmt.Printf("route request failed, err %v", m.err)
		response.Err = m.err
		m.r <- response
		return
	}

	response.Err = json.NewDecoder(resp.Body).Decode(&response.Resp)
	m.r <- response

}

func (oc *OsrmClient) handle() {
	routeChan := make(chan message)
	//tableChan := make(chan message)

	for {
		select {
		case routeReq := <-oc.routeRequest:
			go oc.handleRoute(routeReq, routeChan)
		case m := <-routeChan:
			go handleRouteResponse(&m)
		}
	}
}
