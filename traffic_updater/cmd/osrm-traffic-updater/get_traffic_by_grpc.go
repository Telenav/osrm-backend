package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	proxy "github.com/Telenav/osrm-backend/traffic_updater/pkg/gen-trafficproxy"
	"google.golang.org/grpc"
)

const (
	proxyConnectionTimeout = 60 * time.Second
	maxMsgSize             = 1024 * 1024 * 1024
)

func quickViewFlows(flows []*proxy.FlowResponse, viewCount int) {
	for i := 0; i < viewCount && i < len(flows); i++ {
		log.Printf("--->quickViewFlows %d: %v\n", i, flows[i])
	}
}
func quickViewIncidents(incidents []*proxy.IncidentResponse, viewCount int) {
	for i := 0; i < viewCount && i < len(incidents); i++ {
		log.Printf("--->quickViewIncidents %d: %v\n", i, incidents[i])
	}
}

func trafficData2map(trafficData proxy.TrafficResponse, m map[int64]int) {
	startTime := time.Now()
	defer func() {
		log.Printf("Processing time for building traffic map takes %f seconds\n", time.Now().Sub(startTime).Seconds())
	}()

	var fwdCnt, bwdCnt uint64
	for _, flow := range trafficData.FlowResponses {
		wayid := flow.Flow.WayId
		m[wayid] = int(flow.Flow.Speed)

		if wayid > 0 {
			fwdCnt++
		} else {
			bwdCnt++
		}
	}

	//TODO: support incidents

	log.Printf("Load map[wayid] to speed with %d items, %d forward and %d backward.\n", (fwdCnt + bwdCnt), fwdCnt, bwdCnt)
}

func getTrafficFlowsIncidentsByGRPC(f trafficProxyFlags, wayIds []int64) (*proxy.TrafficResponse, error) {
	var outTrafficResponse proxy.TrafficResponse

	startTime := time.Now()
	defer func() {
		log.Printf("Processing time for getting traffic flows,incidents(%d,%d) for %d wayIds takes %f seconds\n",
			len(outTrafficResponse.FlowResponses), len(outTrafficResponse.IncidentResponses), len(wayIds), time.Now().Sub(startTime).Seconds())
	}()

	// make RPC client
	targetServer := f.ip + ":" + strconv.Itoa(f.port)
	log.Println("dialing traffic proxy " + targetServer)
	conn, err := grpc.Dial(targetServer, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, fmt.Errorf("fail to dial: %v", err)
	}
	defer conn.Close()

	// prepare context
	ctx, cancel := context.WithTimeout(context.Background(), proxyConnectionTimeout)
	defer cancel()

	// new proxy client
	client := proxy.NewTrafficServiceClient(conn)

	// get flows
	var req proxy.TrafficRequest
	req.TrafficSource = new(proxy.TrafficSource)
	req.TrafficSource.Region = f.region
	req.TrafficSource.TrafficProvider = f.trafficProvider
	req.TrafficSource.MapProvider = f.mapProvider
	req.TrafficType = append(req.TrafficType, proxy.TrafficType_FLOW, proxy.TrafficType_INCIDENT)
	if len(wayIds) > 0 {
		log.Printf("getting flows,incidents for %d wayIds\n", len(wayIds))
		var trafficWayIdsRequest proxy.TrafficRequest_TrafficWayIdsRequest
		trafficWayIdsRequest.TrafficWayIdsRequest = new(proxy.TrafficWayIdsRequest)
		trafficWayIdsRequest.TrafficWayIdsRequest.WayIds = wayIds
		req.RequestOneof = &trafficWayIdsRequest
	} else {
		log.Printf("getting all flows,incidents\n")
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

func getDeltaTrafficFlowsIncidentsByGRPCStreaming(f trafficProxyFlags, out chan<- proxy.TrafficResponse) error {
	defer close(out)

	// make RPC client
	targetServer := f.ip + ":" + strconv.Itoa(f.port)
	log.Println("connect traffic proxy " + targetServer)
	conn, err := grpc.Dial(targetServer, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return fmt.Errorf("fail to dial: %v", err)
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
	req.TrafficSource.Region = f.region
	req.TrafficSource.TrafficProvider = f.trafficProvider
	req.TrafficSource.MapProvider = f.mapProvider
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
