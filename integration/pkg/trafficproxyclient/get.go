package trafficproxyclient

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// GetFlowsIncidents return flows and incidents for wayIds or full region.
func GetFlowsIncidents(wayIDs []int64) (*trafficproxy.TrafficResponse, error) {
	var outTrafficResponse trafficproxy.TrafficResponse
	forStr := "all"
	if len(wayIDs) > 0 {
		forStr = fmt.Sprintf("%d wayIds", len(wayIDs))
	}

	startTime := time.Now()
	defer func() {
		glog.Infof("Processing time for getting traffic flows,incidents(%d,%d) for %s takes %f seconds\n",
			len(outTrafficResponse.FlowResponses), len(outTrafficResponse.IncidentResponses),
			forStr, time.Now().Sub(startTime).Seconds())
	}()

	// make RPC client
	conn, err := newGRPCConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// prepare context
	ctx, cancel := context.WithTimeout(context.Background(), params{}.rpcGetTimeout())
	defer cancel()

	// new trafficproxy client
	client := trafficproxy.NewTrafficServiceClient(conn)

	// get flows
	glog.Infof("getting flows,incidents for %s\n", forStr)
	var req trafficproxy.TrafficRequest
	req.TrafficSource = params{}.newTrafficSource()
	req.TrafficType = params{}.newTrafficType()
	if len(wayIDs) > 0 {
		var trafficWayIDsRequest trafficproxy.TrafficRequest_TrafficWayIDsRequest
		trafficWayIDsRequest.TrafficWayIDsRequest = new(trafficproxy.TrafficWayIDsRequest)
		trafficWayIDsRequest.TrafficWayIDsRequest.WayIDs = wayIDs
		req.RequestOneof = &trafficWayIDsRequest
	} else {
		trafficAllRequest := new(trafficproxy.TrafficRequest_TrafficAllRequest)
		trafficAllRequest.TrafficAllRequest = new(trafficproxy.TrafficAllRequest)
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
