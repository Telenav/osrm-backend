package ranking

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/Telenav/osrm-backend/integration/pkg/backend"
	"github.com/golang/glog"
)

// New create a new single host reverseproxy with set timeout for backend OSRM.
func New(target string) *httputil.ReverseProxy {

	u := url.URL{
		Scheme: "http",
		Host:   target,
	}
	targetProxy := httputil.NewSingleHostReverseProxy(&u)
	targetProxy.Transport = backend.NewTransport()
	targetProxy.ErrorHandler = backend.ErrorHandleFunc
	targetProxy.ModifyResponse = modifyResponse

	return targetProxy
}

func modifyResponse(resp *http.Response) error {
	glog.V(3).Info(resp)
	glog.Infof("osrm reverse proxy access,  Request: %v Accept-Encoding: %s, Status: %s ContentLength: %d Content-Encoding: %s",
		resp.Request.URL, resp.Request.Header.Get("Accept-Encoding"), resp.Status, resp.ContentLength, resp.Header.Get("Content-Encoding"))
	return nil
}
