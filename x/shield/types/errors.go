package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrNotShieldOperator = sdkerrors.Register(ModuleName, 101, "not the shield operator account")
	ErrNoDeposit         = sdkerrors.Register(ModuleName, 102, "no coins given for initial deposit")
	ErrNoShield          = sdkerrors.Register(ModuleName, 103, "no coins given for shield")
	ErrEmptySponsor      = sdkerrors.Register(ModuleName, 104, "no sponsor specified for a pool")
	ErrNoPoolFound       = sdkerrors.Register(ModuleName, 105, "no pool found for the given ID")
	ErrNoUpdate          = sdkerrors.Register(ModuleName, 106, "nothing was updated for the pool")
	ErrInvalidGenesis    = sdkerrors.Register(ModuleName, 107, "invalid genesis state")
	ErrInvalidPoolID     = sdkerrors.Register(ModuleName, 108, "invalid pool ID")
	ErrInvalidDuration   = sdkerrors.Register(ModuleName, 109, "invalid specification of coverage duration")
	ErrCannotExtend      = sdkerrors.Register(ModuleName, 110,
		"invalid type (time in seconds or number of blocks) specified to extend this pool")
	ErrNoDelegationAmount  = sdkerrors.Register(ModuleName, 111, "cannot obtain delegation amount info")
	ErrInsufficientStaking = sdkerrors.Register(ModuleName, 112,
		"insufficient total delegation amount to deposit the collateral")
	ErrPoolAlreadyPaused          = sdkerrors.Register(ModuleName, 113, "pool is already paused")
	ErrPoolAlreadyActive          = sdkerrors.Register(ModuleName, 114, "pool is already active")
	ErrPoolInactive               = sdkerrors.Register(ModuleName, 115, "pool is inactive")
	ErrPurchaseMissingDescription = sdkerrors.Register(ModuleName, 116, "missing description for the purchase")
	ErrNotEnoughShield            = sdkerrors.Register(ModuleName, 117, "not enough available shield")
	ErrNoPurchaseFound            = sdkerrors.Register(ModuleName, 118, "no purchase found for the given txhash")
	ErrNoRewards                  = sdkerrors.Register(ModuleName, 119, "no foreign coins rewards to transfer for the denomination")
	ErrInvalidDenom               = sdkerrors.Register(ModuleName, 120, "invalid coin denomination")
	ErrInvalidToAddr              = sdkerrors.Register(ModuleName, 121, "invalid recipient address")
	ErrEmptySender                = sdkerrors.Register(ModuleName, 122, "no sender provided")
	ErrPoolLifeTooShort           = sdkerrors.Register(ModuleName, 123, "new pool life is too short")
)
