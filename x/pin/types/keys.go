package types

const (
	// ModuleName defines the module name
	ModuleName = "pin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_pin"
)

var (
	ParamsKey = []byte("p_pin")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
