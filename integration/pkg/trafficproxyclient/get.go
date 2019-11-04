package trafficproxyclient

import (
	"context"
	"fmt"
	"io"
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
	"github.com/golang/glog"
)

// GetFlowsIncidents return flows and incidents for wayIds or full region.
func GetFlowsIncidents(wayIds []int64) (*proxy.TrafficResponse, error) {
	var outTrafficResponse proxy.TrafficResponse
	forStr := "all"
	if len(wayIds) > 0 {
		forStr = fmt.Sprintf("%d wayIds", len(wayIds))
	}

	startTime := time.Now()
	defer func() {
		glog.Infof("Processing time for getting traffic flows,incidents(%d,%d) for %s takes %f seconds\n",
			len(outTrafficResponse.FlowResponses), len(outTrafficResponse.IncidentResponses),
			forStr, time.Now().Sub(startTime).Seconds())
	}()

	// make RPC client
	conn, err := NewGRPCConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// prepare context
	ctx, cancel := context.WithTimeout(context.Background(), params{}.rpcGetTimeout())
	defer cancel()

	// new proxy client
	client := proxy.NewTrafficServiceClient(conn)

	// get flows
	glog.Infof("getting flows,incidents for %s\n", forStr)
	var req proxy.TrafficRequest
	req.TrafficSource = params{}.newTrafficSource()
	req.TrafficType = params{}.newTrafficType()
	if len(wayIds) > 0 {
		var trafficWayIdsRequest proxy.TrafficRequest_TrafficWayIdsRequest
		trafficWayIdsRequest.TrafficWayIdsRequest = new(proxy.TrafficWayIdsRequest)
		trafficWayIdsRequest.TrafficWayIdsRequest.WayIds = wayIds
		req.RequestOneof = &trafficWayIdsRequest
	} else {
		trafficAllRequest := new(proxy.TrafficRequest_TrafficAllRequest)
		trafficAllRequest.TrafficAllRequest = new(proxy.TrafficAllRequest)
		req.RequestOneof = trafficAllRequest
	}

	glog.V(2).Infof("rpc request: %v", req)
	stream, err := client.GetTrafficData(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("GetTrafficData failed, err: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stream recv failed, err: %v", err)
		}
		glog.V(2).Infof("received traffic data from stream, got flows count: %d, incidents count: %d\n", len(resp.FlowResponses), len(resp.IncidentResponses))
		outTrafficResponse.FlowResponses = append(outTrafficResponse.FlowResponses, resp.FlowResponses...)
		outTrafficResponse.IncidentResponses = append(outTrafficResponse.IncidentResponses, resp.IncidentResponses...)
	}

	return &outTrafficResponse, nil
}
