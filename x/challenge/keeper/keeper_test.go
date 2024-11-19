package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/stretchr/testify/require"

	topTypes "topchain/types"
	"topchain/x/challenge/keeper"
	challenge "topchain/x/challenge/module"
	"topchain/x/challenge/types"
)

const (
	Alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	Bob   = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	Carol = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
)

// Custom BankKeeper interface because MintCoins is not to be defined in the main one.
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddress sdk.AccAddress, recipientModule string, amount sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddress sdk.AccAddress, amount sdk.Coins) error
	// Methods imported from bank should be defined here
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
}

func MockChallengeKeeper(t testing.TB) (keeper.Keeper, sdk.Context, challenge.AppModule) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	logger := log.NewNopLogger()

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	storeService := runtime.NewKVStoreService(storeKey)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)

	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	module := authtypes.NewModuleAddress(types.ModuleName)
	// stakingAccount := authtypes.NewModuleAddress(stakingtypes.ModuleName)

	var maccPerms = map[string][]string{
		authtypes.FeeCollectorName: nil,
		// ... other module accounts ...
		types.ModuleName:    {authtypes.Minter, authtypes.Burner},
		govtypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	}

	accountKeeper := authkeeper.NewAccountKeeper(cdc, storeService, authtypes.ProtoBaseAccount,
		maccPerms, authcodec.NewBech32Codec("cosmos"), "cosmos", string(authority))

	bankKeeper := bankkeeper.NewBaseKeeper(cdc, storeService, accountKeeper, nil, authority.String(), logger)

	k := keeper.NewKeeper(
		cdc,
		storeService,
		logger,
		authority.String(),
		module.String(),
		accountKeeper,
		bankKeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, logger)
	// Prefund accounts
	if err := MockFundAccounts(bankKeeper, ctx); err != nil {
		panic(err)
	}
	appModule := challenge.NewAppModule(cdc, k, accountKeeper, bankKeeper)
	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx, appModule
}

func MockBlockHeight(ctx sdk.Context, am challenge.AppModule, height int64) sdk.Context {
	header := cmtproto.Header{Height: height}
	ctx = sdk.NewContext(ctx.MultiStore(), header, false, log.NewNopLogger())
	_ = am.EndBlock(ctx)
	return ctx
}

func MockFundAccounts(bankKeeper BankKeeper, ctx sdk.Context) error {
	totalMint := sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, 10000000000))
	amounts := sdk.NewCoins(sdk.NewInt64Coin(topTypes.TokenDenom, 1000000000))
	if err := bankKeeper.MintCoins(ctx, types.ModuleName, totalMint); err != nil {
		return err
	}

	for _, addr := range []string{Alice, Bob, Carol} {
		recipientAddress, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return err
		}
		if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddress, amounts); err != nil {
			return err
		}
	}
	return nil
}
