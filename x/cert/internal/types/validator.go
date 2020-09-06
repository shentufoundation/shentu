package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validator is a type for certified validator.
type Validator struct {
	PubKey    crypto.PubKey
	Certifier sdk.AccAddress
}

// Validators is a collection of Validator objects.
type Validators []Validator
