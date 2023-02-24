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
	errProgramNotAllowed
	errProgramExpired
	errProgramPubKey
	errProgramID
)

// Finding
const (
	errFindingNotExists uint32 = iota + 201
	errFindingStatusInvalid
	errFindingSubmitterInvalid
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
	ErrProgramNotExists            = sdkerrors.Register(ModuleName, errProgramNotExists, "program does not exist")
	ErrProgramInactive             = sdkerrors.Register(ModuleName, errProgramInactive, "program is inactive")
	ErrProgramCreatorInvalid       = sdkerrors.Register(ModuleName, errProgramCreatorInvalid, "invalid program creator")
	ErrProgramNotAllowed           = sdkerrors.Register(ModuleName, errProgramNotAllowed, "program access denied because you are not the creator or certifiers")
	ErrProgramExpired              = sdkerrors.Register(ModuleName, errProgramExpired, "cannot end an expired program")
	ErrProgramPubKey               = sdkerrors.Register(ModuleName, errProgramPubKey, "invalid program public key")
	ErrProgramID                   = sdkerrors.Register(ModuleName, errProgramID, "invalid program id")
)

// [2xx] Finding
var (
	ErrFindingNotExists            = sdkerrors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid        = sdkerrors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingSubmitterInvalid     = sdkerrors.Register(ModuleName, errFindingSubmitterInvalid, "invalid finding submitter")
	ErrFindingPlainTextDataInvalid = sdkerrors.Register(ModuleName, errFindingPlainTextDataInvalid, "invalid finding plain text data")
	ErrFindingEncryptedDataInvalid = sdkerrors.Register(ModuleName, errFindingEncryptedDataInvalid, "invalid finding encrypted data")
	ErrFindingID                   = sdkerrors.Register(ModuleName, errFindingID, "invalid finding id")
)

var ErrInvalidGenesis = sdkerrors.Register(ModuleName, errInvalidGenesis, "invalid genesis state")
