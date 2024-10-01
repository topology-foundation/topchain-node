package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"topchain/x/subscription/types"

	query "github.com/cosmos/cosmos-sdk/types/query"

	keepertest "topchain/testutil/keeper"
)

func TestSubscription(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	subscription := types.Subscription{
		Id:       "sub1",
		Provider: "provider1",
	}

	keeper.SetSubscription(ctx, subscription)

	req := &types.QuerySubscriptionRequest{Id: "sub1"}

	retrievedSubscription, err := keeper.Subscription(ctx, req)
	require.NoError(t, err)
	require.Equal(t, subscription, retrievedSubscription.Subscription)
}

func TestSubscriptions(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	subscription1 := types.Subscription{
		Id:       "sub1",
		Provider: "provider1",
	}
	subscription2 := types.Subscription{
		Id:       "sub2",
		Provider: "provider1",
	}
	subscription3 := types.Subscription{
		Id:       "sub3",
		Provider: "provider2",
	}

	keeper.SetSubscription(ctx, subscription1)
	keeper.SetSubscription(ctx, subscription2)
	keeper.SetSubscription(ctx, subscription3)

	req := &types.QuerySubscriptionsRequest{Provider: "provider1"}
	res, err := keeper.Subscriptions(ctx, req)
	require.NoError(t, err)
	require.Len(t, res.Subscriptions, 2)
	require.Contains(t, res.Subscriptions, subscription1)
	require.Contains(t, res.Subscriptions, subscription2)
}

func TestSubscriptionsWithPaginationOne(t *testing.T) {
	keeper, ctx := keepertest.SubscriptionKeeper(t)
	subscription1 := types.Subscription{
		Id:       "sub1",
		Provider: "provider1",
	}
	subscription2 := types.Subscription{
		Id:       "sub2",
		Provider: "provider1",
	}
	subscription3 := types.Subscription{
		Id:       "sub3",
		Provider: "provider2",
	}

	keeper.SetSubscription(ctx, subscription1)
	keeper.SetSubscription(ctx, subscription2)
	keeper.SetSubscription(ctx, subscription3)

	req := &types.QuerySubscriptionsRequest{Provider: "provider1", Pagination: &query.PageRequest{Limit: 1}}
	res, err := keeper.Subscriptions(ctx, req)
	require.NoError(t, err)
	require.Len(t, res.Subscriptions, 1)
	require.Contains(t, res.Subscriptions, subscription1)
	req = &types.QuerySubscriptionsRequest{Provider: "provider1", Pagination: &query.PageRequest{Key: res.Pagination.NextKey, Limit: 1}}
	res, err = keeper.Subscriptions(ctx, req)
	require.Len(t, res.Subscriptions, 1)
	require.Contains(t, res.Subscriptions, subscription2)
}
