package searchconnector

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Telenav/osrm-backend/integration/api/search/nearbychargestation"
	"github.com/Telenav/osrm-backend/integration/pkg/backend"
	"github.com/golang/glog"
)

type tnSearchHTTPClient struct {
	tnSearchEndpoint string
	apiKey           string
	apiSignature     string
	httpclient       http.Client
	requestC         chan *request
}

func newTNSearchHTTPClient(endpoint, apiKey, apiSignature string) *tnSearchHTTPClient {
	return &tnSearchHTTPClient{
		tnSearchEndpoint: endpoint,
		apiKey:           apiKey,
		apiSignature:     apiSignature,
		httpclient:       http.Client{Timeout: backend.Timeout()},
		requestC:         make(chan *request),
	}
}

func (sc *tnSearchHTTPClient) submitSearchReq(req *nearbychargestation.Request) <-chan ChargeStationsResponse {
	var url string
	if !strings.HasPrefix(sc.tnSearchEndpoint, "http://") {
		url += "http://"
	}
	req.APIKey = sc.apiKey
	req.APISignature = sc.apiSignature
	url = url + sc.tnSearchEndpoint + req.RequestURI()

	searchReq := newTNSearchRequest(url)
	sc.requestC <- searchReq
	return searchReq.searchRespC
}

func (sc *tnSearchHTTPClient) start() {
	glog.Info("search connector started.\n")
	for {
		select {
		case req := <-sc.requestC:
			go sc.handle(req)
		}
	}
}

func (sc *tnSearchHTTPClient) handle(req *request) {
	defer close(req.searchRespC)

	resp, err := sc.httpclient.Get(req.url)
	glog.Infof("[tnSearchHTTPClient]Finish http request for %s" + req.url)
	if err != nil || resp == nil {
		glog.Warningf("search request %s failed, err %v\n", req.url, err)
	}
	defer resp.Body.Close()

	var searchResp ChargeStationsResponse
	searchResp.Err = json.NewDecoder(resp.Body).Decode(&searchResp.Resp)
	req.searchRespC <- searchResp
	glog.Infof("[tnSearchHTTPClient]Response for request %s" + req.url + "is generated.")
}
