package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Program
const (
	errProgramFindingListEmpty uint32 = iota + 101
	errProgramFindingListMarshal
	errProgramFindingListUnmarshal
	errProgramAlreadyExists
	errProgramNotExists
	errProgramAlreadyActive
	errProgramAlreadyClosed
	errProgramInactive
	errProgramNotInactive
	errProgramStatusInvalid
	errProgramCreatorInvalid
	errProgramNotAllowed
	errProgramExpired
	errProgramID
)

// Finding
const (
	errFindingAlreadyExists uint32 = iota + 201
	errFindingNotExists
	errFindingStatusInvalid
	errFindingHashInvalid
	errFindingSeverityLevelInvalid
	errFindingSubmitterInvalid
	errFindingNotAllowed
	errFindingPlainTextDataInvalid
	errFindingEncryptedDataInvalid
	errFindingID
)

const errInvalidGenesis = 301

// [1xx] Program
var (
	ErrProgramFindingListEmpty     = sdkerrors.Register(ModuleName, errProgramFindingListEmpty, "empty finding id list")
	ErrProgramFindingListMarshal   = sdkerrors.Register(ModuleName, errProgramFindingListMarshal, "convert uint64 to byte list error")
	ErrProgramFindingListUnmarshal = sdkerrors.Register(ModuleName, errProgramFindingListUnmarshal, "convert to uint64 list error")
	ErrProgramAlreadyExists        = sdkerrors.Register(ModuleName, errProgramAlreadyExists, "program already exists")
	ErrProgramNotExists            = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exists")
	ErrProgramAlreadyActive        = sdkerrors.Register(ModuleName, errProgramAlreadyActive, "program already active")
	ErrProgramAlreadyClosed        = sdkerrors.Register(ModuleName, errProgramAlreadyClosed, "program already closed")
	ErrProgramNotActive            = sdkerrors.Register(ModuleName, errProgramInactive, "program status is not active")
	ErrProgramNotInactive          = sdkerrors.Register(ModuleName, errProgramNotInactive, "program status is not inactive")
	ErrProgramStatusInvalid        = sdkerrors.Register(ModuleName, errProgramStatusInvalid, "program status invalid")
	ErrProgramCreatorInvalid       = sdkerrors.Register(ModuleName, errProgramCreatorInvalid, "invalid program creator")
	ErrProgramOperatorNotAllowed   = sdkerrors.Register(ModuleName, errProgramNotAllowed, "program access denied")
	ErrProgramExpired              = sdkerrors.Register(ModuleName, errProgramExpired, "cannot end an expired program")
	ErrProgramID                   = sdkerrors.Register(ModuleName, errProgramID, "invalid program id")
)

// [2xx] Finding
var (
	ErrFindingAlreadyExists        = sdkerrors.Register(ModuleName, errFindingAlreadyExists, "program already exists")
	ErrFindingNotExists            = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid        = sdkerrors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingHashInvalid          = sdkerrors.Register(ModuleName, errFindingHashInvalid, "invalid finding hash")
	ErrFindingSeverityLevelInvalid = sdkerrors.Register(ModuleName, errFindingSeverityLevelInvalid, "invalid finding severity level")

	ErrFindingSubmitterInvalid     = sdkerrors.Register(ModuleName, errFindingSubmitterInvalid, "invalid finding submitter")
	ErrFindingOperatorNotAllowed   = sdkerrors.Register(ModuleName, errFindingNotAllowed, "finding access denied")
	ErrFindingPlainTextDataInvalid = sdkerrors.Register(ModuleName, errFindingPlainTextDataInvalid, "invalid finding plain text data")
	ErrFindingEncryptedDataInvalid = sdkerrors.Register(ModuleName, errFindingEncryptedDataInvalid, "invalid finding encrypted data")
	ErrFindingID                   = sdkerrors.Register(ModuleName, errFindingID, "invalid finding id")
)

var ErrInvalidGenesis = sdkerrors.Register(ModuleName, errInvalidGenesis, "invalid genesis state")
