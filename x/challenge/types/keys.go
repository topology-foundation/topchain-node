package types

const (
	// ModuleName defines the module name
	ModuleName = "challenge"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_challenge"

	ChallengeKeyPrefix = "Challenge/value"
	ProofKeyPrefix     = "Proof/value"
)

var (
	ParamsKey = []byte("p_challenge")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetProofStoreKey(challengeId string) []byte {
	return KeyPrefix(ProofKeyPrefix + "/" + challengeId)
}
