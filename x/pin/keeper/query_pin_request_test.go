package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "topchain/testutil/keeper"
	"topchain/testutil/nullify"
	"topchain/x/pin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestPinRequestQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.PinKeeper(t)
	msgs := createNPinRequest(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetPinRequestRequest
		response *types.QueryGetPinRequestResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPinRequestRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetPinRequestResponse{PinRequest: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPinRequestRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetPinRequestResponse{PinRequest: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPinRequestRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.PinRequest(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestPinRequestQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.PinKeeper(t)
	msgs := createNPinRequest(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPinRequestRequest {
		return &types.QueryAllPinRequestRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PinRequestAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PinRequest), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PinRequest),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PinRequestAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PinRequest), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PinRequest),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PinRequestAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.PinRequest),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.PinRequestAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
