package osrmconnector

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Telenav/osrm-backend/integration/api/osrm/route"
	"github.com/Telenav/osrm-backend/integration/api/osrm/table"
)

var fakeOSRMRouteResponse = route.Response{
	Code:    "OK",
	Message: "RouteResponse",
}

var fakeOSRMTableResponse = table.Response{
	Code:    "OK",
	Message: "TableResponse",
}

func TestSingleRouteRequest(t *testing.T) {
	var osrmRouteResponseBytes, _ = json.Marshal(fakeOSRMRouteResponse)
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
	var osrmRouteResponseBytes, _ = json.Marshal(fakeOSRMRouteResponse)
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
			count++
		}
		if count == 10 {
			close(all)
		}
	}

	if count != 10 {
		t.Errorf("TestMultiRouteRequest failed!")
	}
}

func TestMixRouteTableRequest(t *testing.T) {
	var osrmRouteResponseBytes, _ = json.Marshal(fakeOSRMRouteResponse)
	var osrmTableResponseBytes, _ = json.Marshal(fakeOSRMTableResponse)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}

		if r.URL.EscapedPath() == "/route/v1/driving/" {
			w.Write(osrmRouteResponseBytes)
		} else {
			w.Write(osrmTableResponseBytes)
		}
	}))
	defer ts.Close()

	oc := newOsrmHTTPClient(ts.URL)
	go oc.start()

	clientCount := 20
	all := make(chan bool, clientCount)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < clientCount; i++ {
		if rand.Intn(2) == 0 {
			// send route request and verify
			go func() {
				req := route.NewRequest()
				result := oc.submitRouteReq(req)
				r := <-result
				if r.Err == nil && r.Resp.Code == "OK" && r.Resp.Message == "RouteResponse" {
					all <- true
				} else {
					all <- false
				}
			}()
		} else {
			// send table request and verify
			go func() {
				req := table.NewRequest()
				result := oc.submitTableReq(req)
				r := <-result
				if r.Err == nil && r.Resp.Code == "OK" && r.Resp.Message == "TableResponse" {
					all <- true
				} else {
					all <- false
				}
			}()
		}
	}

	go func() {
		time.Sleep(5 * time.Second)
		close(all)
	}()

	count := 0
	for b := range all {
		if b {
			count++
		}

		if count == clientCount {
			close(all)
		}
	}

	if count != clientCount {
		t.Errorf("TestMixRouteTableRequest failed!")
	}
}
