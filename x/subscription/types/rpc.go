package types

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client TopologyRpcClient

func SetupClient() (TopologyRpcClient, error) {
	// TODO: open issue default, we need to allow custom configs
	conn, err := grpc.NewClient("localhost:6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client = NewTopologyRpcClient(conn)
	return client, nil
}

func GetClient() TopologyRpcClient {
	return client
}
