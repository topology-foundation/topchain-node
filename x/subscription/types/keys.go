package types

const (
	// ModuleName defines the module name
	ModuleName = "subscription"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscription"

	DealKeyPrefix                 = "Deal/value"
	SubscriptionKeyPrefix         = "Subscription/value"
	SubscriptionProviderKeyPrefix = "Subscription/Provider/value"
)

var ParamsKey = []byte("p_subscription")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
