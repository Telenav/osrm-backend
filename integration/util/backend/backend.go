package backend

import (
	"flag"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
)

var flags struct {
	timeout time.Duration
}

func init() {
	flag.DurationVar(&flags.timeout, "backend-timeout", 60000000000, "Timeout for sending request and waiting response against backend endpoint.") // default 60 seconds
}

// Timeout returns expect timeout for sending request and waiting response against backend endpoint.
func Timeout() time.Duration {
	return flags.timeout
}

// NewTransport create a new transport for backend services.
func NewTransport() *http.Transport {

	// refer to default transport, remove proxy and changed timeout
	// https://golang.org/pkg/net/http/#RoundTripper
	t := http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   Timeout(),
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &t
}

// ErrorHandleFunc provides a general way to handle error for backend services.
func ErrorHandleFunc(w http.ResponseWriter, req *http.Request, err error) {
	glog.Errorf("backend request %v error: %v", req.URL, err)
	http.Error(w, err.Error(), http.StatusBadGateway)
}
