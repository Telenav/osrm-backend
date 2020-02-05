package osrmconnector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrm/route"
)

var fakeOsrmRouteResponse = route.Response{
	Code: "OK",
}

func TestSingleRouteRequest(t *testing.T) {
	var osrmRouteResponseBytes, _ = json.Marshal(fakeOsrmRouteResponse)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(osrmRouteResponseBytes)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
	}))
	defer ts.Close()

	oc := newOsrmHTTPClient(ts.URL)
	go oc.start()
	req := route.NewRequest()
	result := oc.submitRouteReq(req)
	c := <-result
	jsonResult, _ := json.Marshal(c.Resp)
	fmt.Println(string(jsonResult))
}

func TestMultiRouteRequest(t *testing.T) {
	var osrmRouteResponseBytes, _ = json.Marshal(fakeOsrmRouteResponse)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(osrmRouteResponseBytes)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
	}))
	defer ts.Close()

	oc := newOsrmHTTPClient(ts.URL)
	go oc.start()

	all := make(chan RouteResponse, 10)

	for i := 0; i < 10; i++ {
		go func() {
			req := route.NewRequest()
			result := oc.submitRouteReq(req)
			all <- (<-result)
		}()
	}

	go func() {
		time.Sleep(3 * time.Second)
		close(all)
	}()

	count := 0
	for r := range all {
		if r.Resp.Code == "OK" {
			count += 1
		}
		if count == 10 {
			close(all)
		}
	}

	if count != 10 {
		t.Errorf("TestMultiRouteRequest failed!")
	}
}
