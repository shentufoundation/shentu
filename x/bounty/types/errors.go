package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Error Code Enums

const (
	errUnknownHost uint32 = iota + 101
	errUnknownProgram
	errUnknownFinding
)

// x/bounty module sentinel errors
var (
	ErrUnknownHost    = sdkerrors.Register(ModuleName, errUnknownHost, "unknown host")
	ErrUnknownProgram = sdkerrors.Register(ModuleName, errUnknownProgram, "unknown program")
	ErrUnknownFinding = sdkerrors.Register(ModuleName, errUnknownFinding, "unknown finding")
)
