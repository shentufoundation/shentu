package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal = sdkerrors.Register(ModuleName, 2, "unknown proposal")
	ErrInvalidGenesis  = sdkerrors.Register(ModuleName, 3, "invalid genesis state")
)
