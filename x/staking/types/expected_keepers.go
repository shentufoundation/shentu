package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CertKeeper interface {
	GetValidator(ctx sdk.Context, validator cryptotypes.PubKey) ([]byte, bool)
	IsValidatorCertified(ctx sdk.Context, validator cryptotypes.PubKey) bool
}
