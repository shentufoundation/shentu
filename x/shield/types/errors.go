package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotShieldOperator = sdkerrors.Register(ModuleName, 1, "not the shield admin account")
	ErrNoDeposit         = sdkerrors.Register(ModuleName, 2, "no coins given for initial deposit")
	ErrNoShield          = sdkerrors.Register(ModuleName, 3, "no coins given for shield shield")
	ErrEmptySponsor      = sdkerrors.Register(ModuleName, 4, "no sponsor specified for a pool")
	ErrNoPoolFound       = sdkerrors.Register(ModuleName, 5, "no pool found for the given ID")
	ErrNoUpdate          = sdkerrors.Register(ModuleName, 6, "nothing was updated for the pool")
	ErrInvalidGenesis    = sdkerrors.Register(ModuleName, 7, "invalid genesis state")
	ErrInvalidPoolID     = sdkerrors.Register(ModuleName, 8, "invalid pool ID")
	ErrInvalidDuration   = sdkerrors.Register(ModuleName, 9, "invalid specification of coverage duration")
	ErrCannotExtend      = sdkerrors.Register(ModuleName, 10,
		"invalid type (time in seconds or number of blocks) specified to extend this pool")

	ErrPoolAlreadyPaused = sdkerrors.Register(ModuleName, 100, "pool is already paused")
	ErrPoolAlreadyActive = sdkerrors.Register(ModuleName, 101, "pool is already active")
)
