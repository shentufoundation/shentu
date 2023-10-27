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
	errProgramAlreadyActive
	errProgramAlreadyClosed
	errProgramNotExists
	errProgramInactive
	errProgramNotInactive
	errProgramCreatorInvalid
	errProgramNotAllowed
	errProgramExpired
	errProgramPubKey
	errProgramID
	errNoProgramFound
)

// Finding
const (
	errFindingAlreadyExists uint32 = iota + 201
	errFindingNotExists
	errFindingStatusInvalid
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
	ErrProgramAlreadyActive        = sdkerrors.Register(ModuleName, errProgramAlreadyActive, "program already active")
	ErrProgramAlreadyClosed        = sdkerrors.Register(ModuleName, errProgramAlreadyClosed, "program already closed")

	ErrProgramNotExists   = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exists")
	ErrProgramNotActive   = sdkerrors.Register(ModuleName, errProgramInactive, "program status is not active")
	ErrProgramNotInactive = sdkerrors.Register(ModuleName, errProgramNotInactive, "program status is not inactive")

	ErrProgramCreatorInvalid     = sdkerrors.Register(ModuleName, errProgramCreatorInvalid, "invalid program creator")
	ErrProgramOperatorNotAllowed = sdkerrors.Register(ModuleName, errProgramNotAllowed, "program access denied because you are not the creator or certifiers")
	ErrProgramExpired            = sdkerrors.Register(ModuleName, errProgramExpired, "cannot end an expired program")
	ErrProgramPubKey             = sdkerrors.Register(ModuleName, errProgramPubKey, "invalid program public key")
	ErrProgramID                 = sdkerrors.Register(ModuleName, errProgramID, "invalid program id")
	ErrNoProgramFound            = sdkerrors.Register(ModuleName, errNoProgramFound, "program does not exist")
)

// [2xx] Finding
var (
	ErrFindingAlreadyExists        = sdkerrors.Register(ModuleName, errFindingAlreadyExists, "program already exists")
	ErrFindingNotExists            = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid        = sdkerrors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingSubmitterInvalid     = sdkerrors.Register(ModuleName, errFindingSubmitterInvalid, "invalid finding submitter")
	ErrFindingOperatorNotAllowed   = sdkerrors.Register(ModuleName, errFindingNotAllowed, "finding access denied because you are not the creator or certifiers")
	ErrFindingPlainTextDataInvalid = sdkerrors.Register(ModuleName, errFindingPlainTextDataInvalid, "invalid finding plain text data")
	ErrFindingEncryptedDataInvalid = sdkerrors.Register(ModuleName, errFindingEncryptedDataInvalid, "invalid finding encrypted data")
	ErrFindingID                   = sdkerrors.Register(ModuleName, errFindingID, "invalid finding id")
)

var ErrInvalidGenesis = sdkerrors.Register(ModuleName, errInvalidGenesis, "invalid genesis state")
