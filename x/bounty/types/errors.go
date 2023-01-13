package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// Error Code Enums
// Program
const (
	errProgramNotExists uint32 = iota + 101
)

// Finding
const (
	errFindingNotExists uint32 = iota + 201
	errFindingInvalid
	errFindingAccessDeny
	errFindingAlreadyInactive
	errFindingAlreadyActive
)

// [1xx] Program
var (
	ErrProgramNotExists = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exist")
)

// [2xx] Finding
var (
	ErrFindingNotExists       = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingInvalid         = sdkerrors.Register(ModuleName, errFindingInvalid, "invalid finding content")
	ErrFindingAccessDeny      = sdkerrors.Register(ModuleName, errFindingAccessDeny, "not the finding submitter")
	ErrFindingAlreadyInactive = sdkerrors.Register(ModuleName, errFindingAlreadyInactive, "finding already inactive")
	ErrFindingAlreadyActive   = sdkerrors.Register(ModuleName, errFindingAlreadyActive, "finding already active")
)
