package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Program
const (
	errProgramAlreadyExists uint32 = iota + 101
	errProgramNotExists
	errProgramAlreadyActive
	errProgramAlreadyClosed
	errProgramInactive
	errProgramStatusInvalid
	errProgramOperatorNotAllowed
	errProgramCloseNotAllowed
	errProgramID
)

// Finding
const (
	errFindingAlreadyExists uint32 = iota + 201
	errFindingNotExists
	errFindingStatusInvalid
	errFindingHashInvalid
	errFindingSeverityLevelInvalid
	errFindingOperatorNotAllowed
	errFindingID
)

// [1xx] Program
var (
	ErrProgramAlreadyExists      = sdkerrors.Register(ModuleName, errProgramAlreadyExists, "program already exists")
	ErrProgramNotExists          = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exists")
	ErrProgramAlreadyActive      = sdkerrors.Register(ModuleName, errProgramAlreadyActive, "program already active")
	ErrProgramAlreadyClosed      = sdkerrors.Register(ModuleName, errProgramAlreadyClosed, "program already closed")
	ErrProgramNotActive          = sdkerrors.Register(ModuleName, errProgramInactive, "program status is not active")
	ErrProgramStatusInvalid      = sdkerrors.Register(ModuleName, errProgramStatusInvalid, "program status invalid")
	ErrProgramOperatorNotAllowed = sdkerrors.Register(ModuleName, errProgramOperatorNotAllowed, "program access denied")
	ErrProgramCloseNotAllowed    = sdkerrors.Register(ModuleName, errProgramCloseNotAllowed, "cannot close the program")
	ErrProgramID                 = sdkerrors.Register(ModuleName, errProgramID, "invalid program id")
)

// [2xx] Finding
var (
	ErrFindingAlreadyExists        = sdkerrors.Register(ModuleName, errFindingAlreadyExists, "finding already exists")
	ErrFindingNotExists            = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid        = sdkerrors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingHashInvalid          = sdkerrors.Register(ModuleName, errFindingHashInvalid, "invalid finding hash")
	ErrFindingSeverityLevelInvalid = sdkerrors.Register(ModuleName, errFindingSeverityLevelInvalid, "invalid finding severity level")
	ErrFindingOperatorNotAllowed   = sdkerrors.Register(ModuleName, errFindingOperatorNotAllowed, "finding access denied")
	ErrFindingID                   = sdkerrors.Register(ModuleName, errFindingID, "invalid finding id")
)
