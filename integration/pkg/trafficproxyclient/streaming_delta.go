package trafficproxyclient

import (
	"context"
	"fmt"
	"io"

	"github.com/Telenav/osrm-backend/integration/pkg/trafficproxy"
	"github.com/golang/glog"
)

// StreamingDeltaFlowsIncidents set up a new channel for traffic flows and incidents streaming delta.
func StreamingDeltaFlowsIncidents(out chan<- trafficproxy.TrafficResponse) error {

	// make RPC client
	conn, err := newGRPCConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	// prepare context
	ctx := context.Background()

	// new proxy client
	client := trafficproxy.NewTrafficServiceClient(conn)

	// get flows via stream
	glog.Info("streaming delta traffic flows,incidents")
	var req trafficproxy.TrafficRequest
	req.TrafficSource = params{}.newTrafficSource()
	req.TrafficType = params{}.newTrafficType()
	trafficDeltaStreamRequest := new(trafficproxy.TrafficRequest_TrafficStreamingDeltaRequest)
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest = new(trafficproxy.TrafficStreamingDeltaRequest)
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest.StreamingRule = params{}.newStreamingRule()
	req.RequestOneof = trafficDeltaStreamRequest

	glog.V(2).Infof("rpc request: %v", req)
	stream, err := client.GetTrafficData(ctx, &req)
	if err != nil {
		return fmt.Errorf("GetTrafficData failed, err: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			glog.Info("streaming delta has finished gracefully.")
			break
		}
		if err != nil {
			return fmt.Errorf("stream recv failed, err: %v", err)
		}
		glog.V(2).Infof("received traffic data from stream, got flows count: %d, incidents count: %d\n", len(resp.FlowResponses), len(resp.IncidentResponses))
		out <- *resp
	}

	return nil
}
