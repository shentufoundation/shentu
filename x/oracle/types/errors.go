package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
)

const errInconsistentOperators uint32 = 301

var (
	ErrNoOperatorFound         = sdkerrors.Register(ModuleName, errNoOperatorFound, "no operator was found")
	ErrOperatorAlreadyExists   = sdkerrors.Register(ModuleName, errOperatorAlreadyExists, "operator already exists")
	ErrInvalidDueBlock         = sdkerrors.Register(ModuleName, errInvalidDueBlock, "invalid due block")
	ErrNoTotalCollateralFound  = sdkerrors.Register(ModuleName, errNoTotalCollateralFound, "total collateral not found")
	ErrNoEnoughTotalCollateral = sdkerrors.Register(ModuleName, errNoEnoughTotalCollateral, "total collateral not enough")
	ErrTotalCollateralNotEqual = sdkerrors.Register(ModuleName, errTotalCollateralNotEqual, "total collateral not equal")
	ErrNoEnoughCollateral      = sdkerrors.Register(ModuleName, errNoEnoughCollateral, "collateral not enough")
	ErrInvalidPoolParams       = sdkerrors.Register(ModuleName, errInvalidPoolParams, "invalid pool params")
	ErrInvalidTaskParams       = sdkerrors.Register(ModuleName, errInvalidTaskParams, "invalid task params")

	ErrTaskNotExists       = sdkerrors.Register(ModuleName, errTaskNotExists, "task does not exist")
	ErrUnqualifiedOperator = sdkerrors.Register(ModuleName, errUnqualifiedOperator, "operator is not qualified")
	ErrDuplicateResponse   = sdkerrors.Register(ModuleName, errDuplicateResponse, "already receive response from this operator")
	ErrTaskClosed          = sdkerrors.Register(ModuleName, errTaskClosed, "task is already closed")
	ErrTaskNotClosed       = sdkerrors.Register(ModuleName, errTaskNotClosed, "task has not been closed")
	ErrNotExpired          = sdkerrors.Register(ModuleName, errNotExpired, "task is not expired")
	ErrNotCreator          = sdkerrors.Register(ModuleName, errNotCreated, "only creator is allowed to perform this action")
	ErrNotFinished         = sdkerrors.Register(ModuleName, errNotFinished, "the task is on going")
	ErrTaskFailed          = sdkerrors.Register(ModuleName, errTaskFailed, "task failed")
	ErrInvalidScore        = sdkerrors.Register(ModuleName, errInvalidScore, "invalid score")

	ErrInconsistentOperators = sdkerrors.Register(ModuleName, errInconsistentOperators, "two operators not consistent")
)
