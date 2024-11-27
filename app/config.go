package app

import (
	mandutypes "mandu/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	// Set prefixes
	accountPubKeyPrefix := mandutypes.AccountAddressPrefix + "pub"
	validatorAddressPrefix := mandutypes.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := mandutypes.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := mandutypes.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := mandutypes.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(mandutypes.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
