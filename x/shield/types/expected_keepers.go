package types

import (
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ParamSubspace defines the expected Subspace interface for parameters (noalias)
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	// IterateValidators iterates through validators by admin address, execute func for each validator.
	IterateValidators(sdk.Context, func(index int64, validator stakingtypes.ValidatorI) (stop bool))

	// GetValidator gets a particular validator by admin address with a found flag.
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, bool)
	// GetAllValidators gets the set of all validators with no limits, used during genesis dump.
	GetAllValidators(ctx sdk.Context) []stakingtypes.Validator
	// GetValidatorDelegations returns all delegations to a specific validator. Useful for querier.
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) []stakingtypes.Delegation

	// Delegation allows for getting a particular delegation for a given validator
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
	GetAllUnbondingDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.UnbondingDelegation
	GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.UnbondingDelegation, bool)
	SetUnbondingDelegation(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation)
	RemoveUnbondingDelegation(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation)
	GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (dvPairs []stakingtypes.DVPair)
	SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, timeslice []stakingtypes.DVPair)
	InsertUBDQueue(ctx sdk.Context, ubd stakingtypes.UnbondingDelegation, completionTime time.Time)
	SetDelegation(ctx sdk.Context, delegation stakingtypes.Delegation)
	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, bool)
	UBDQueueIterator(ctx sdk.Context, timestamp time.Time) sdk.Iterator
	RemoveValidatorTokensAndShares(ctx sdk.Context, validator stakingtypes.Validator, sharesToRemove sdk.Dec) (valOut stakingtypes.Validator, removedTokens math.Int)
	RemoveUBDQueue(ctx sdk.Context, timestamp time.Time)
	GetRedelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (redelegations []stakingtypes.Redelegation)
	SetValidator(ctx sdk.Context, validator stakingtypes.Validator)
	DeleteValidatorByPowerIndex(ctx sdk.Context, validator stakingtypes.Validator)
	RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) error
	RemoveValidator(ctx sdk.Context, address sdk.ValAddress)

	BondDenom(sdk.Context) string
	UnbondingTime(sdk.Context) time.Duration
	GetBondedPool(ctx sdk.Context) authtypes.ModuleAccountI
}

// BankKeeper defines the expected bank keeper.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoins(ctx sdk.Context, fromAdd, toAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// GovKeeper defines the expected gov keeper.
type GovKeeper interface {
	GetVotingParams(ctx sdk.Context) govtypesv1.VotingParams
}
