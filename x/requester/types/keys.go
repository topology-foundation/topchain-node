package types

const (
	// ModuleName defines the module name
	ModuleName = "requester"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_requester"

	SubscriptionKey = "Subscription/value"
)

var (
	ParamsKey = []byte("p_requester")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
