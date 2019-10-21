package trafficproxyclient

import (
	"context"
	"fmt"
	"io"
	"log"

	proxy "github.com/Telenav/osrm-backend/traffic_updater/pkg/gen-trafficproxy"
)

// NewStreamingDeltaFlowsIncidents set up a new channel for traffic flows and incidents streaming delta.
func NewStreamingDeltaFlowsIncidents(out chan<- proxy.TrafficResponse) error {
	defer close(out)

	// make RPC client
	conn, err := NewGRPCConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	// prepare context
	ctx := context.Background()

	// new proxy client
	client := proxy.NewTrafficServiceClient(conn)

	// get flows via stream
	log.Println("getting delta traffic flows,incidents via stream")
	var req proxy.TrafficRequest
	req.TrafficSource = new(proxy.TrafficSource)
	req.TrafficSource.Region = flags.Region
	req.TrafficSource.TrafficProvider = flags.TrafficProvider
	req.TrafficSource.MapProvider = flags.MapProvider
	req.TrafficType = append(req.TrafficType, proxy.TrafficType_FLOW, proxy.TrafficType_INCIDENT)
	trafficDeltaStreamRequest := new(proxy.TrafficRequest_TrafficStreamingDeltaRequest)
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest = new(proxy.TrafficStreamingDeltaRequest)
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest.StreamingRule = new(proxy.TrafficStreamingDeltaRequest_StreamingRule)
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest.StreamingRule.MaxSize = 1000
	trafficDeltaStreamRequest.TrafficStreamingDeltaRequest.StreamingRule.MaxTime = 5
	req.RequestOneof = trafficDeltaStreamRequest

	stream, err := client.GetTrafficData(ctx, &req)
	if err != nil {
		return fmt.Errorf("GetTrafficData failed, err: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("stream recv failed, err: %v", err)
		}
		log.Printf("[VERBOSE] received traffic data from stream, got flows count: %d, incidents count: %d\n", len(resp.FlowResponses), len(resp.IncidentResponses))
		out <- *resp
	}

	return nil
}
