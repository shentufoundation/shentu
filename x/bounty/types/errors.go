package types

import (
	"cosmossdk.io/errors"
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
	ErrProgramAlreadyExists      = errors.Register(ModuleName, errProgramAlreadyExists, "program already exists")
	ErrProgramNotExists          = errors.Register(ModuleName, errProgramNotExists, "program does not exists")
	ErrProgramAlreadyActive      = errors.Register(ModuleName, errProgramAlreadyActive, "program already active")
	ErrProgramAlreadyClosed      = errors.Register(ModuleName, errProgramAlreadyClosed, "program already closed")
	ErrProgramNotActive          = errors.Register(ModuleName, errProgramInactive, "program status is not active")
	ErrProgramStatusInvalid      = errors.Register(ModuleName, errProgramStatusInvalid, "program status invalid")
	ErrProgramOperatorNotAllowed = errors.Register(ModuleName, errProgramOperatorNotAllowed, "program access denied")
	ErrProgramCloseNotAllowed    = errors.Register(ModuleName, errProgramCloseNotAllowed, "cannot close the program")
	ErrProgramID                 = errors.Register(ModuleName, errProgramID, "invalid program id")
)

// [2xx] Finding
var (
	ErrFindingAlreadyExists        = errors.Register(ModuleName, errFindingAlreadyExists, "finding already exists")
	ErrFindingNotExists            = errors.Register(ModuleName, errFindingNotExists, "finding does not exist")
	ErrFindingStatusInvalid        = errors.Register(ModuleName, errFindingStatusInvalid, "invalid finding status")
	ErrFindingHashInvalid          = errors.Register(ModuleName, errFindingHashInvalid, "invalid finding hash")
	ErrFindingSeverityLevelInvalid = errors.Register(ModuleName, errFindingSeverityLevelInvalid, "invalid finding severity level")
	ErrFindingOperatorNotAllowed   = errors.Register(ModuleName, errFindingOperatorNotAllowed, "finding access denied")
	ErrFindingID                   = errors.Register(ModuleName, errFindingID, "invalid finding id")
)
