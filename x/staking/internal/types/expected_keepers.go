package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type CertKeeper interface {
	GetValidator(ctx sdk.Context, validator crypto.PubKey) ([]byte, bool)
	IsValidatorCertified(ctx sdk.Context, validator crypto.PubKey) bool
}
