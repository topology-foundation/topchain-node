package keeper

import (
	"context"
	"topchain/x/subscription/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Deal(goCtx context.Context, req *types.QueryDealRequest) (*types.QueryDealResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryDealResponse{Deal: deal}, nil
}

func (k Keeper) DealStatus(goCtx context.Context, req *types.QueryDealStatusRequest) (*types.QueryDealStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	deal, found := k.GetDeal(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryDealStatusResponse{Status: deal.Status}, nil
}

func (k Keeper) Deals(ctx context.Context, req *types.QueryDealsRequest) (*types.QueryDealsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DealKeyPrefix))

	var deals []types.Deal
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var deal types.Deal
		if err := k.cdc.Unmarshal(value, &deal); err != nil {
			return err
		}

		if deal.Requester == req.Requester {
			deals = append(deals, deal)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDealsResponse{Deals: deals, Pagination: pageRes}, nil
}
