package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
)

type (
	AccountKeeper interface {
		GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	}

	StakingKeeper interface {
		ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI
	}

	SlashingKeeper interface {
		IsTombstoned(sdk.Context, sdk.ConsAddress) bool
		Tombstone(sdk.Context, sdk.ConsAddress)
		Jail(sdk.Context, sdk.ConsAddress)
		JailUntil(sdk.Context, sdk.ConsAddress, time.Time)
	}
)
