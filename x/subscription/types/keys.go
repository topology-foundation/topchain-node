package types

const (
	// ModuleName defines the module name
	ModuleName = "subscription"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscription"

	DealKeyPrefix                 = "Deal/value"
	DealRequesterKeyPrefix        = "Deal/Requester/value"
	SubscriptionKeyPrefix         = "Subscription/value"
	SubscriptionProviderKeyPrefix = "Subscription/Provider/value"
)

var ParamsKey = []byte("p_subscription")

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetProviderStoreKey returns the key for the provider store for the given provider.
func GetProviderStoreKey(provider string) []byte {
	return []byte(KeyPrefix(SubscriptionProviderKeyPrefix + "/" + provider))
}

// GetRequesterStoreKey returns the key for the requester store for the given requester.
func GetRequesterStoreKey(requester string) []byte {
	return []byte(KeyPrefix(DealRequesterKeyPrefix + "/" + requester))
}
