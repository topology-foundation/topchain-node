package types

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client TopologyRpcClient

func SetupClient() (TopologyRpcClient, error) {
	conn, err := grpc.NewClient("localhost:6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	Client = NewTopologyRpcClient(conn)
	return Client, nil
}
