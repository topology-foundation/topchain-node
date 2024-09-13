package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PinRequestList: []PinRequest{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in pinRequest
	pinRequestIndexMap := make(map[string]struct{})

	for _, elem := range gs.PinRequestList {
		index := string(PinRequestKey(elem.Index))
		if _, ok := pinRequestIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for pinRequest")
		}
		pinRequestIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
