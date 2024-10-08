package types

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var RpcClient TopologyRpcClient

func SetupRpcClient() (TopologyRpcClient, error) {
	conn, err := grpc.NewClient("localhost:6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	RpcClient = NewTopologyRpcClient(conn)
	return RpcClient, nil
}
