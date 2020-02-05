package osrmconnector

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/pkg/backend"
	"github.com/golang/glog"
)

type osrmHTTPClient struct {
	osrmBackendEndpoint string
	httpclient          http.Client
	requestC            chan *request
}

func newOsrmHTTPClient(osrmEndpoint string) *osrmHTTPClient {
	return &osrmHTTPClient{
		osrmBackendEndpoint: osrmEndpoint,
		httpclient:          http.Client{Timeout: backend.Timeout()},
		requestC:            make(chan *request),
	}
}

func (oc *osrmHTTPClient) submitRouteReq(r *route.Request) <-chan RouteResponse {
	var url string
	if !strings.HasPrefix(oc.osrmBackendEndpoint, "http://") {
		url += "http://"
	}
	url = url + oc.osrmBackendEndpoint + r.RequestURI()

	req := newOsrmRequest(url, OSRMROUTE)
	oc.requestC <- req
	return req.routeRespC
}

func (oc *osrmHTTPClient) start() {
	c := make(chan message)

	for {
		select {
		case req := <-oc.requestC:
			go oc.send(req, c)
		case m := <-c:
			go oc.response(&m)
		}
	}
}

type message struct {
	req  *request
	resp *http.Response
	err  error
}

func (oc *osrmHTTPClient) send(req *request, c chan<- message) {
	resp, err := oc.httpclient.Get(req.url)
	m := message{req: req, resp: resp, err: err}
	c <- m
}

func (oc *osrmHTTPClient) response(m *message) {
	defer close(m.req.routeRespC)
	defer close(m.req.tableRespC)

	var routeResp RouteResponse
	if m.err != nil || m.resp == nil {
		glog.Warningf("osrm request %s failed, err %v\n", m.req.url, m.err)

		if m.req.t == OSRMROUTE {
			routeResp.Err = m.err
			m.req.routeRespC <- routeResp
		}

		return
	}
	defer m.resp.Body.Close()

	if m.req.t == OSRMROUTE {
		routeResp.Err = json.NewDecoder(m.resp.Body).Decode(&routeResp.Resp)
		m.req.routeRespC <- routeResp
	}

	return
}
