package rankingservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Telenav/osrm-backend/integration/pkg/api/osrmv1"
	"github.com/golang/glog"
)

func (h *Handler) routeByOSRM(osrmRequest *osrmv1.RouteRequest) (*osrmv1.RouteResponse, int, error) {
	if osrmRequest == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("empty osrm request")
	}

	osrmRequestURL := "http://" + h.osrmBackend + osrmRequest.RequestURI()
	glog.Infof("osrm request to backend: %s", osrmRequestURL)

	clt := http.Client{Timeout: h.backendTimeout}
	resp, err := clt.Get(osrmRequestURL)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("route request %s agianst OSRM failed, err %v", osrmRequestURL, err)
	}
	defer resp.Body.Close()

	// parse OSRM response
	var routeResponse osrmv1.RouteResponse
	err = json.NewDecoder(resp.Body).Decode(&routeResponse)
	if err != nil {
		glog.Warningf("Decode osrm HTTP response body failed, http status %d, err %v.", resp.StatusCode, err)
		glog.V(3).Info(resp.Body)
		return nil, resp.StatusCode, err
	}
	glog.Infof("osrm response from backend, http status %d, response code %s, message %s, data_version %s",
		resp.StatusCode, routeResponse.Code, routeResponse.Message, routeResponse.DataVersion)
	glog.V(3).Infof("osrm response from backend: %v", routeResponse)

	return &routeResponse, resp.StatusCode, nil
}
