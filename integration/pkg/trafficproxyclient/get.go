package trafficproxyclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	proxy "github.com/Telenav/osrm-backend/integration/pkg/gen-trafficproxy"
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
		log.Printf("Processing time for getting traffic flows,incidents(%d,%d) for %s takes %f seconds\n",
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// new proxy client
	client := proxy.NewTrafficServiceClient(conn)

	// get flows
	log.Printf("getting flows,incidents for %s\n", forStr)
	var req proxy.TrafficRequest
	req.TrafficSource = new(proxy.TrafficSource)
	req.TrafficSource.Region = flags.Region
	req.TrafficSource.TrafficProvider = flags.TrafficProvider
	req.TrafficSource.MapProvider = flags.MapProvider
	req.TrafficType = append(req.TrafficType, proxy.TrafficType_FLOW, proxy.TrafficType_INCIDENT)
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
		log.Printf("[VERBOSE] received traffic data from stream, got flows count: %d, incidents count: %d\n", len(resp.FlowResponses), len(resp.IncidentResponses))
		outTrafficResponse.FlowResponses = append(outTrafficResponse.FlowResponses, resp.FlowResponses...)
		outTrafficResponse.IncidentResponses = append(outTrafficResponse.IncidentResponses, resp.IncidentResponses...)
	}

	return &outTrafficResponse, nil
}
