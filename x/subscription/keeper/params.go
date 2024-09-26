package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"topchain/x/subscription/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx context.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k Keeper) CreateModuleAccount(ctx sdk.Context, moduleName string) {

	// Check if the module account already exists
	macc := k.accountKeeper.GetModuleAccount(ctx, moduleName)
	if macc == nil {
		// 	return nil // Module account already exists, no action needed
		// moduleAcc := authtypes.NewEmptyModuleAccount(moduleName, authtypes.Minter, authtypes.Burner)
		// k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
		// fmt.Println("Created new module account: %s\n", moduleName)

	}
}

func (k Keeper) GetModuleAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleAccountName)
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleAccountName)
}
