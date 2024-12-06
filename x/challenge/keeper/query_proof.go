package keeper

import (
	"context"

	"mandu/x/challenge/types"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Proof(goCtx context.Context, req *types.QueryProofRequest) (*types.QueryProofResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	vertex, found := k.GetProof(ctx, req.ChallengeId, req.Hash)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryProofResponse{Vertex: vertex}, nil
}

func (k Keeper) Proofs(goCtx context.Context, req *types.QueryProofsRequest) (*types.QueryProofsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.GetProofStoreKey(req.ChallengeId))

	var vertices []types.Vertex
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, _ []byte) error {
		vertex, found := k.GetProof(ctx, req.ChallengeId, string(key))
		if !found {
			return sdkerrors.ErrKeyNotFound
		}

		vertices = append(vertices, vertex)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProofsResponse{Vertices: vertices, Pagination: pageRes}, nil
}
