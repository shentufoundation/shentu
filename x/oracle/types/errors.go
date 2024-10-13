package types

import (
	"cosmossdk.io/errors"
)

// Error Code Enums
const (
	errNoOperatorFound uint32 = iota + 101
	errOperatorAlreadyExists
	errInvalidDueBlock
	errNoTotalCollateralFound
	errNoEnoughTotalCollateral
	errTotalCollateralNotEqual
	errNoEnoughCollateral
	errInvalidPoolParams
	errInvalidTaskParams
	errUnqualifiedRemover
	errUnqualifiedCreator
)

const (
	errTaskNotExists uint32 = iota + 201
	errUnqualifiedOperator
	errDuplicateResponse
	errTaskClosed
	errTaskNotClosed
	errNotExpired
	errNotCreated
	errNotFinished
	errTaskFailed
	errInvalidScore
	errInvalidTask
	errOverdueValidTime
	errUnexpectedTask
	errTooLateValidTime
)

const errInconsistentOperators uint32 = 301
const errFailedToCastTask uint32 = 401

var (
	ErrNoOperatorFound         = errors.Register(ModuleName, errNoOperatorFound, "no operator was found")
	ErrOperatorAlreadyExists   = errors.Register(ModuleName, errOperatorAlreadyExists, "operator already exists")
	ErrInvalidDueBlock         = errors.Register(ModuleName, errInvalidDueBlock, "invalid due block")
	ErrNoTotalCollateralFound  = errors.Register(ModuleName, errNoTotalCollateralFound, "total collateral not found")
	ErrNoEnoughTotalCollateral = errors.Register(ModuleName, errNoEnoughTotalCollateral, "total collateral not enough")
	ErrTotalCollateralNotEqual = errors.Register(ModuleName, errTotalCollateralNotEqual, "total collateral not equal")
	ErrNoEnoughCollateral      = errors.Register(ModuleName, errNoEnoughCollateral, "collateral not enough")
	ErrInvalidPoolParams       = errors.Register(ModuleName, errInvalidPoolParams, "invalid pool params")
	ErrInvalidTaskParams       = errors.Register(ModuleName, errInvalidTaskParams, "invalid task params")
	ErrUnqualifiedRemover      = errors.Register(ModuleName, errUnqualifiedRemover, "unauthorized to remove this operator - not the operator itself nor a certifier")
	ErrUnqualifiedCreator      = errors.Register(ModuleName, errUnqualifiedCreator, "unauthorized to create this operator - not a certifier")

	ErrTaskNotExists       = errors.Register(ModuleName, errTaskNotExists, "task does not exist")
	ErrUnqualifiedOperator = errors.Register(ModuleName, errUnqualifiedOperator, "operator is not qualified")
	ErrDuplicateResponse   = errors.Register(ModuleName, errDuplicateResponse, "already receive response from this operator")
	ErrTaskClosed          = errors.Register(ModuleName, errTaskClosed, "task is already closed")
	ErrTaskNotClosed       = errors.Register(ModuleName, errTaskNotClosed, "task has not been closed")
	ErrNotExpired          = errors.Register(ModuleName, errNotExpired, "task is not expired")
	ErrNotCreator          = errors.Register(ModuleName, errNotCreated, "only creator is allowed to perform this action")
	ErrNotFinished         = errors.Register(ModuleName, errNotFinished, "the task is on going")
	ErrTaskFailed          = errors.Register(ModuleName, errTaskFailed, "task failed")
	ErrInvalidScore        = errors.Register(ModuleName, errInvalidScore, "invalid score")
	ErrInvalidTask         = errors.Register(ModuleName, errInvalidTask, "invalid task")
	ErrOverdueValidTime    = errors.Register(ModuleName, errOverdueValidTime, "the valid time is overdue")
	ErrUnexpectedTask      = errors.Register(ModuleName, errUnexpectedTask, "a different typed task already exists")
	ErrTooLateValidTime    = errors.Register(ModuleName, errTooLateValidTime, "the valid time is later than expiration time")

	ErrInconsistentOperators = errors.Register(ModuleName, errInconsistentOperators, "two operators not consistent")

	ErrFailedToCastTask = errors.Register(ModuleName, errFailedToCastTask, "failed to cast to concrete task")
)
