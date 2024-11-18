package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"topchain/x/challenge/keeper"
	challenge "topchain/x/challenge/module"
	"topchain/x/challenge/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, sdk.Context, challenge.AppModule) {
	k, ctx, am := MockChallengeKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx, am
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx, am := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotNil(t, am)
	require.NotEmpty(t, k)
}
