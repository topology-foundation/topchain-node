package requester

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "topchain/api/topchain/requester"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ListSubscription",
					Use:            "list-subscription",
					Short:          "Query list-subscription",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "RequestSubscription",
					Use:            "request-subscription [cro_id] [amount] [duration]",
					Short:          "Send a request-subscription tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "cro_id"}, {ProtoField: "amount"}, {ProtoField: "duration"}},
				},
				{
					RpcMethod:      "CancelSubscription",
					Use:            "cancel-subscription [subscription_id]",
					Short:          "Send a cancel-subscription tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "subscription_id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
