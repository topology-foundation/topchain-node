package subscription

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"topchain/testutil/sample"
	subscriptionsimulation "topchain/x/subscription/simulation"
	"topchain/x/subscription/types"
)

// avoid unused import issue
var (
	_ = subscriptionsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgRequestSubscription = "op_weight_msg_request_subscription"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRequestSubscription int = 100

	opWeightMsgCancelSubscription = "op_weight_msg_cancel_subscription"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCancelSubscription int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	subscriptionGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&subscriptionGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgRequestSubscription int
	simState.AppParams.GetOrGenerate(opWeightMsgRequestSubscription, &weightMsgRequestSubscription, nil,
		func(_ *rand.Rand) {
			weightMsgRequestSubscription = defaultWeightMsgRequestSubscription
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRequestSubscription,
		subscriptionsimulation.SimulateMsgRequestSubscription(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCancelSubscription int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelSubscription, &weightMsgCancelSubscription, nil,
		func(_ *rand.Rand) {
			weightMsgCancelSubscription = defaultWeightMsgCancelSubscription
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelSubscription,
		subscriptionsimulation.SimulateMsgCancelSubscription(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgRequestSubscription,
			defaultWeightMsgRequestSubscription,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				subscriptionsimulation.SimulateMsgRequestSubscription(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCancelSubscription,
			defaultWeightMsgCancelSubscription,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				subscriptionsimulation.SimulateMsgCancelSubscription(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
