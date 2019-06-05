package main

import (
	"context"
	"fmt"

	"github.com/Telenav/osrm-backend/traffic_updater/gen-go/proxy"
	"github.com/apache/thrift/lib/go/thrift"
)

var defaultCtx = context.Background()

func main() {

	var transport thrift.TTransport
	var err error

	fmt.Println("connect traffic proxy ")
	transport, err = thrift.NewTSocket("127.0.0.1:6666")
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return
	}

	transport, err = thrift.NewTFramedTransportFactoryMaxLength(thrift.NewTTransportFactory(), 1024*1024*1024).GetTransport(transport)
	if err != nil {
		fmt.Println("Error get transport:", err)
		return
	}
	defer transport.Close()

	if err := transport.Open(); err != nil {
		return
	}
	iprot := thrift.NewTCompactProtocolFactory().GetProtocol(transport)
	oprot := thrift.NewTCompactProtocolFactory().GetProtocol(transport)
	client := proxy.NewProxyServiceClient(thrift.NewTStandardClient(iprot, oprot))

	fmt.Println("getting flows")
	flows, err := client.GetAllFlows(defaultCtx)
	if err != nil {
		fmt.Println("get flows failed:", err)
		return
	}
	fmt.Printf("got flows count: %d\n", len(flows))
}
