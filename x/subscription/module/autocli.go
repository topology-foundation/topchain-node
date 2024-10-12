package subscription

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "topchain/api/topchain/subscription"
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
					RpcMethod:      "Deal",
					Use:            "deal [id]",
					Short:          "Query deal",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "DealStatus",
					Use:            "deal-status [id]",
					Short:          "Query deal status",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "Deals",
					Use:            "deals [requester]",
					Short:          "Query deals",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "requester"}},
				},
				{
					RpcMethod:      "Subscription",
					Use:            "subscription [id]",
					Short:          "Query subscription",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "Subscriptions",
					Use:            "subscriptions [provider]",
					Short:          "Query subscriptions",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider"}},
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
					RpcMethod:      "CreateDeal",
					Use:            "create-deal [cro_id] [amount] [start_block] [end_block] [initial_frontier]",
					Short:          "Send a create-deal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "cro_id"}, {ProtoField: "amount"}, {ProtoField: "start_block"}, {ProtoField: "end_block"}, {ProtoField: "initial_frontier"}},
				},
				{
					RpcMethod:      "CancelDeal",
					Use:            "cancel-deal [deal_id]",
					Short:          "Send a cancel-deal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deal_id"}},
				},
				{
					RpcMethod:      "UpdateDeal",
					Use:            "update-deal [deal_id] [amount] [start_block] [end_block]",
					Short:          "Send a update-deal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deal_id"}, {ProtoField: "amount"}, {ProtoField: "start_block"}, {ProtoField: "end_block"}},
				},
				{
					RpcMethod:      "IncrementDealAmount",
					Use:            "increment-deal-amount [deal_id] [amount]",
					Short:          "Send a increment-deal-amount tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deal_id"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "JoinDeal",
					Use:            "join-deal [deal_id]",
					Short:          "Send a join-deal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deal_id"}},
				},
				{
					RpcMethod:      "LeaveDeal",
					Use:            "leave-deal [deal_id]",
					Short:          "Send a leave-deal tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "deal_id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
