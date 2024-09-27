package types

const (
	// ModuleName defines the module name
	ModuleName = "subscription"

	ModuleAccountName = "subscription_account"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscription"

	DealKeyPrefix         = "Deal/value"
	SubscriptionKeyPrefix = "Subscription/value"
)

var ParamsKey = []byte("p_subscription")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
