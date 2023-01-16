package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Error Code Enums
const (
	errUnknownProgram uint32 = iota + 101
	errUnknownHost
	errUnknownFinding
)

const (
	ErrorEmptyProgramIDFindingList = "empty finding id list"
)

// x/bounty module sentinel errors
var (
	ErrUnknownProgram = sdkerrors.Register(ModuleName, errUnknownProgram, "unknown program")
	ErrUnknownHost    = sdkerrors.Register(ModuleName, errUnknownHost, "unknown host")
	ErrUnknownFinding = sdkerrors.Register(ModuleName, errUnknownFinding, "unknown finding")
)
