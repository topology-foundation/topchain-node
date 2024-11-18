package challenge

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "topchain/api/topchain/challenge"
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
					RpcMethod:      "Proof",
					Use:            "proof [challenge_id] [hash]",
					Short:          "Query proof",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}, {ProtoField: "hash"}},
				},
				{
					RpcMethod:      "Proofs",
					Use:            "proofs [challenge_id]",
					Short:          "Query proofs",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}},
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
					RpcMethod:      "Challenge",
					Use:            "challenge [provider_id] [vertices_hashes]",
					Short:          "Send a challenge tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "provider_id"}, {ProtoField: "vertices_hashes"}},
				},
				{
					RpcMethod:      "SubmitProof",
					Use:            "submit-proof [challenge_id] [vertices]",
					Short:          "Submit a submit-proof tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}, {ProtoField: "vertices"}},
				},
				{
					RpcMethod:      "RequestDependencies",
					Use:            "request-dependencies [challenge_id] [vertices_hashes]",
					Short:          "Send a request-dependencies tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}, {ProtoField: "vertices_hashes"}},
				},
				{
					RpcMethod:      "SettleChallenge",
					Use:            "settle-challenge [challenge_id]",
					Short:          "Send a settle-challenge tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "challenge_id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
