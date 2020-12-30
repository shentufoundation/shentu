package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoOperatorFound         = sdkerrors.Register(ModuleName, 101, "no operator was found")
	ErrOperatorAlreadyExists   = sdkerrors.Register(ModuleName, 102, "operator already exists")
	ErrInvalidDueBlock         = sdkerrors.Register(ModuleName, 103, "invalid due block")
	ErrNoTotalCollateralFound  = sdkerrors.Register(ModuleName, 104, "total collateral not found")
	ErrNoEnoughTotalCollateral = sdkerrors.Register(ModuleName, 105, "total collateral not enough")
	ErrTotalCollateralNotEqual = sdkerrors.Register(ModuleName, 106, "total collateral not equal")
	ErrNoEnoughCollateral      = sdkerrors.Register(ModuleName, 107, "collateral not enough")
	ErrInvalidPoolParams       = sdkerrors.Register(ModuleName, 108, "invalid pool params")
	ErrInvalidTaskParams       = sdkerrors.Register(ModuleName, 109, "invalid task params")

	ErrTaskNotExists       = sdkerrors.Register(ModuleName, 201, "task does not exist")
	ErrUnqualifiedOperator = sdkerrors.Register(ModuleName, 202, "operator is not qualified")
	ErrDuplicateResponse   = sdkerrors.Register(ModuleName, 203, "already receive response from this operator")
	ErrTaskClosed          = sdkerrors.Register(ModuleName, 204, "task is already closed")
	ErrTaskNotClosed       = sdkerrors.Register(ModuleName, 205, "task has not been closed")
	ErrNotExpired          = sdkerrors.Register(ModuleName, 206, "task is not expired")
	ErrNotCreator          = sdkerrors.Register(ModuleName, 207, "only creator is allowed to perform this action")
	ErrNotFinished         = sdkerrors.Register(ModuleName, 208, "the task is on going")
	ErrTaskFailed          = sdkerrors.Register(ModuleName, 209, "task failed")
	ErrInvalidScore        = sdkerrors.Register(ModuleName, 210, "invalid score")

	ErrInconsistentOperators = sdkerrors.Register(ModuleName, 301, "two operators not consistent")
)
