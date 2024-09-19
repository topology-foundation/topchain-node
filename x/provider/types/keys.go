package types

const (
	// ModuleName defines the module name
	ModuleName = "provider"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_provider"
)

var (
	ParamsKey = []byte("p_provider")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
