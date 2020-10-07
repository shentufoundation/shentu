package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CertKeeper interface {
	GetValidator(ctx sdk.Context, validator crypto.PubKey) ([]byte, bool)
	IsValidatorCertified(ctx sdk.Context, validator crypto.PubKey) bool
}
