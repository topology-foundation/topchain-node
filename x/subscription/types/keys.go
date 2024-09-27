package types

const (
	// ModuleName defines the module name
	ModuleName = "subscription"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscription"

	DealKeyPrefix          = "Deal/value"
	DealRequesterKeyPrefix = "Deal/Requester/value"
	SubscriptionKeyPrefix  = "Subscription/value"
)

var ParamsKey = []byte("p_subscription")

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetRequesterStoreKey returns the key for the provider store for the given provider.
func GetRequesterStoreKey(requester string) []byte {
	return []byte(KeyPrefix(DealRequesterKeyPrefix + "/" + requester))
}
