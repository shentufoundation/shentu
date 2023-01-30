package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Program
const (
	errProgramFindingListEmpty uint32 = iota + 101
	errProgramFindingListMarshal
	errProgramFindingListUnmarshal
	errProgramNotExists
	errProgramInactive
	errProgramCreatorInvalid
)

// Finding
const (
	errFindingNotExists uint32 = iota + 201
	errFindingStatusInvalid
	errFindingSubmitterInvalid
)

// [1xx] Program
var (
	ErrorEmptyProgramIDFindingList   = sdkerrors.Register(ModuleName, errProgramFindingListEmpty, "empty finding id list")
	ErrProgramIDFindingListMarshal   = sdkerrors.Register(ModuleName, errProgramFindingListMarshal, "convert uint64 to byte list error")
	ErrProgramIDFindingListUnmarshal = sdkerrors.Register(ModuleName, errProgramFindingListUnmarshal, "convert to uint64 list error")
	ErrProgramNotExists              = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exist")
	ErrProgramInactive               = sdkerrors.Register(ModuleName, errProgramInactive, "program is inactive")
	ErrProgramCreatorInvalid         = sdkerrors.Register(ModuleName, errProgramCreatorInvalid, "invalid program creator")
)

// [2xx] Finding
var (
	ErrFindingNotExists        = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid    = sdkerrors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingSubmitterInvalid = sdkerrors.Register(ModuleName, errFindingSubmitterInvalid, "invalid finding submitter")
)
