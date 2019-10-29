package trafficproxyclient

import (
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
)

const (
	maxMsgSize = 1024 * 1024 * 1024
)

// NewGRPCConnection create a new GRPC connection to target traffic proxy.
func NewGRPCConnection() (*grpc.ClientConn, error) {

	// make RPC client
	targetServer := flags.IP + ":" + strconv.Itoa(flags.Port)
	log.Println("dialing traffic proxy " + targetServer)
	conn, err := grpc.Dial(targetServer, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, fmt.Errorf("fail to dial traffic proxy %s, err %v", targetServer, err)
	}
	return conn, nil
}
