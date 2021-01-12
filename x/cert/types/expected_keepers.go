package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
)

type (
	AccountKeeper interface {
		GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	}

	StakingKeeper interface {
		ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI
		GetAllValidators(ctx sdk.Context) []staking.Validator
		GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) []staking.Delegation
	}

	SlashingKeeper interface {
		IsTombstoned(sdk.Context, sdk.ConsAddress) bool
		Tombstone(sdk.Context, sdk.ConsAddress)
		Jail(sdk.Context, sdk.ConsAddress)
		JailUntil(sdk.Context, sdk.ConsAddress, time.Time)
	}
)
