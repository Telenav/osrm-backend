package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
	"github.com/apache/thrift/lib/go/thrift"
)

func dumpFlowsToCsv(csv_file string, flows []*proxy.Flow) {

}

func main() {

	var transport thrift.TTransport
	var err error

	// make socket
	fmt.Println("connect traffic proxy ")
	transport, err = thrift.NewTSocket("127.0.0.1:6666")
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return
	}

	// Buffering
	transport, err = thrift.NewTFramedTransportFactoryMaxLength(thrift.NewTTransportFactory(), 1024*1024*1024).GetTransport(transport)
	if err != nil {
		fmt.Println("Error get transport:", err)
		return
	}
	defer transport.Close()
	if err := transport.Open(); err != nil {
		return
	}

	// protocol encoder&decoder
	protocol := thrift.NewTCompactProtocolFactory().GetProtocol(transport)

	// create proxy client
	client := proxy.NewProxyServiceClient(thrift.NewTStandardClient(protocol, protocol))

	// get flows
	startTime := time.Now()
	fmt.Println("getting flows")
	var defaultCtx = context.Background()
	flows, err := client.GetAllFlows(defaultCtx)
	if err != nil {
		fmt.Println("get flows failed:", err)
		return
	}
	fmt.Printf("got flows count: %d\n", len(flows))
	afterGotFlowTime := time.Now()
	fmt.Printf("get flows time used: %f seconds\n", afterGotFlowTime.Sub(startTime).Seconds())

	// TODO: dump to csv
	fmt.Println("dump flows to: ")
	dumpFlowsToCsv("traffic.csv", flows)
	endTime := time.Now()
	fmt.Printf("dump csv time used: %f seconds\n", endTime.Sub(afterGotFlowTime).Seconds())

}
