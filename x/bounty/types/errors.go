package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bounty module sentinel errors
var (
	ErrUnknownProgram = sdkerrors.Register(ModuleName, 2, "unknown program")
	ErrInvalidGenesis = sdkerrors.Register(ModuleName, 3, "invalid genesis state")
)
