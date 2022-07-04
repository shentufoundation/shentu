package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// Error Code Enums

const (
	errNotShieldAdmin uint32 = iota + 101
	errNoDeposit
	errNoShield
	errEmptySponsor
	errNoPoolFound
	errNoUpdate
	errInvalidGenesis
	errInvalidPoolID
	errInvalidDuration
	errAdminWithdraw
	errNoDelegationAmount
	errInsufficientStaking
	errPoolAlreadyPaused
	errPoolAlreadyActive
	errPoolInactive
	errPurchaseMissingDescription
	errNotEnoughShield
	errNoPurchaseFound
	errNoRewards
	errInvalidDenom
	errInvalidToAddr
	errNoCollateralFound
	errInvalidCollateralAmount
	errEmptySender
	errPoolLifeTooShort
	errPurchaseNotFound
	errProviderNotFound
	errNotEnoughCollateral
	errReimbursementNotFound
	errInvalidBeneficiary
	errNotPayoutTime
	errOverWithdraw
	errNoPoolFoundForSponsor
	errSponsorAlreadyExists
	errCollateralBadDenom
	errSponsorPurchase
	errOperationNotSupported
	errPoolShieldExceedsLimit
	errShieldAdminNotActive
	errPurchaseTooSmall
	errNotEnoughStaked
)

var (
	ErrNotShieldAdmin             = sdkerrors.Register(ModuleName, errNotShieldAdmin, "not the shield admin account")
	ErrNoDeposit                  = sdkerrors.Register(ModuleName, errNoDeposit, "no coins given for initial deposit")
	ErrNoShield                   = sdkerrors.Register(ModuleName, errNoShield, "no coins given for shield")
	ErrEmptySponsor               = sdkerrors.Register(ModuleName, errEmptySponsor, "no sponsor specified for a pool")
	ErrNoPoolFound                = sdkerrors.Register(ModuleName, errNoPoolFound, "no pool found")
	ErrNoUpdate                   = sdkerrors.Register(ModuleName, errNoUpdate, "nothing was updated for the pool")
	ErrInvalidGenesis             = sdkerrors.Register(ModuleName, errInvalidGenesis, "invalid genesis state")
	ErrInvalidPoolID              = sdkerrors.Register(ModuleName, errInvalidPoolID, "invalid pool ID")
	ErrInvalidDuration            = sdkerrors.Register(ModuleName, errInvalidDuration, "invalid specification of coverage duration")
	ErrAdminWithdraw              = sdkerrors.Register(ModuleName, errAdminWithdraw, "admin cannot manually withdraw collateral")
	ErrNoDelegationAmount         = sdkerrors.Register(ModuleName, errNoDelegationAmount, "cannot obtain delegation amount info")
	ErrInsufficientStaking        = sdkerrors.Register(ModuleName, errInsufficientStaking, "insufficient total delegation amount to deposit the collateral")
	ErrPoolAlreadyPaused          = sdkerrors.Register(ModuleName, errPoolAlreadyPaused, "pool is already paused")
	ErrPoolAlreadyActive          = sdkerrors.Register(ModuleName, errPoolAlreadyActive, "pool is already active")
	ErrPoolInactive               = sdkerrors.Register(ModuleName, errPoolInactive, "pool is inactive")
	ErrPurchaseMissingDescription = sdkerrors.Register(ModuleName, errPurchaseMissingDescription, "missing description for the purchase")
	ErrNotEnoughShield            = sdkerrors.Register(ModuleName, errNotEnoughShield, "not enough available shield")
	ErrNoPurchaseFound            = sdkerrors.Register(ModuleName, errNoPurchaseFound, "no purchase found for the given txhash")
	ErrNoRewards                  = sdkerrors.Register(ModuleName, errNoRewards, "no foreign coins rewards to transfer for the denomination")
	ErrInvalidDenom               = sdkerrors.Register(ModuleName, errInvalidDenom, "invalid coin denomination")
	ErrInvalidToAddr              = sdkerrors.Register(ModuleName, errInvalidToAddr, "invalid recipient address")
	ErrNoCollateralFound          = sdkerrors.Register(ModuleName, errNoCollateralFound, "no collateral for the pool found with the given provider address")
	ErrInvalidCollateralAmount    = sdkerrors.Register(ModuleName, errInvalidCollateralAmount, "invalid amount of collateral")
	ErrEmptySender                = sdkerrors.Register(ModuleName, errEmptySender, "no sender provided")
	ErrPoolLifeTooShort           = sdkerrors.Register(ModuleName, errPoolLifeTooShort, "new pool life is too short")
	ErrPurchaseNotFound           = sdkerrors.Register(ModuleName, errPurchaseNotFound, "purchase is not found")
	ErrProviderNotFound           = sdkerrors.Register(ModuleName, errProviderNotFound, "provider is not found")
	ErrNotEnoughCollateral        = sdkerrors.Register(ModuleName, errNotEnoughCollateral, "not enough collateral")
	ErrReimbursementNotFound      = sdkerrors.Register(ModuleName, errReimbursementNotFound, "reimbursement is not found")
	ErrInvalidBeneficiary         = sdkerrors.Register(ModuleName, errInvalidBeneficiary, "invalid beneficiary")
	ErrNotPayoutTime              = sdkerrors.Register(ModuleName, errNotPayoutTime, "has not reached payout time yet")
	ErrOverWithdraw               = sdkerrors.Register(ModuleName, errOverWithdraw, "too much withdraw initiated")
	ErrNoPoolFoundForSponsor      = sdkerrors.Register(ModuleName, errNoPoolFoundForSponsor, "no pool found for the given sponsor")
	ErrSponsorAlreadyExists       = sdkerrors.Register(ModuleName, errSponsorAlreadyExists, "a pool already exists under the given sponsor")
	ErrCollateralBadDenom         = sdkerrors.Register(ModuleName, errCollateralBadDenom, "invalid coin denomination for collateral")
	ErrSponsorPurchase            = sdkerrors.Register(ModuleName, errSponsorPurchase, "pool sponsor cannot purchase shield")
	ErrOperationNotSupported      = sdkerrors.Register(ModuleName, errOperationNotSupported, "operation is currently not supported")
	ErrPoolShieldExceedsLimit     = sdkerrors.Register(ModuleName, errPoolShieldExceedsLimit, "pool shield exceeds limit")
	ErrShieldAdminNotActive       = sdkerrors.Register(ModuleName, errShieldAdminNotActive, "shield admin is not activated")
	ErrPurchaseTooSmall           = sdkerrors.Register(ModuleName, errPurchaseTooSmall, "purchase amount is too small")
	ErrNotEnoughStaked            = sdkerrors.Register(ModuleName, errNotEnoughStaked, "not enough unlocked staking to be withdrawn")
)
