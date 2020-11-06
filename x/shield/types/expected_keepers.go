package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	// IterateValidators iterates through validators by admin address, execute func for each validator.
	IterateValidators(sdk.Context, func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// GetValidator gets a particular validator by admin address with a found flag.
	GetValidator(sdk.Context, sdk.ValAddress) (staking.Validator, bool)
	// GetAllValidators gets the set of all validators with no limits, used during genesis dump.
	GetAllValidators(ctx sdk.Context) []staking.Validator
	// GetValidatorDelegations returns all delegations to a specific validator. Useful for querier.
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) []staking.Delegation

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingexported.DelegationI
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []staking.Delegation
	GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) []staking.UnbondingDelegation
	GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (staking.UnbondingDelegation, bool)
	SetUnbondingDelegation(ctx sdk.Context, ubd staking.UnbondingDelegation)
	RemoveUnbondingDelegation(ctx sdk.Context, ubd staking.UnbondingDelegation)
	GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []staking.DVPair)
	SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, timeslice []staking.DVPair)
	InsertUBDQueue(ctx sdk.Context, ubd staking.UnbondingDelegation, completionTime time.Time)
	SetDelegation(ctx sdk.Context, delegation staking.Delegation)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (staking.Delegation, bool)
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
	UBDQueueIterator(ctx sdk.Context, timestamp time.Time) sdk.Iterator
	RemoveValidatorTokensAndShares(ctx sdk.Context, validator staking.Validator, sharesToRemove sdk.Dec) (valOut staking.Validator, removedTokens sdk.Int)
	RemoveUBDQueue(ctx sdk.Context, timestamp time.Time)
	GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []staking.Redelegation)
	SetValidator(ctx sdk.Context, validator staking.Validator)
	DeleteValidatorByPowerIndex(ctx sdk.Context, validator staking.Validator)
	RemoveDelegation(ctx sdk.Context, delegation staking.Delegation)
	RemoveValidator(ctx sdk.Context, address sdk.ValAddress)

	BondDenom(sdk.Context) string
	UnbondingTime(sdk.Context) time.Duration
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoins(ctx sdk.Context, fromAdd, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
}

// SupplyKeeper defines the expected supply keeper.
type SupplyKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// GovKeeper defines the expected gov keeper.
type GovKeeper interface {
	GetVotingParams(ctx sdk.Context) govTypes.VotingParams
}
