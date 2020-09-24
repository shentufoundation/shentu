package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotShieldOperator = sdkerrors.Register(ModuleName, 1, "not the shield operator account")
	ErrNoDeposit         = sdkerrors.Register(ModuleName, 2, "no coins given for initial deposit")
	ErrNoCoverage        = sdkerrors.Register(ModuleName, 3, "no coins given for shield coverage")
	ErrEmptySponsor      = sdkerrors.Register(ModuleName, 4, "no sponsor specified for a pool")
)
