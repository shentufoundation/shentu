package types

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	AccountKeeper interface {
		GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	}

	BankKeeper interface {
		SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	}

	StakingKeeper interface {
		ValidatorByConsAddr(context.Context, sdk.ConsAddress) (stakingtypes.ValidatorI, error) // get a particular validator by consensus address
		GetAllValidators(ctx context.Context) ([]stakingtypes.Validator, error)
		GetValidatorDelegations(ctx context.Context, valAddr sdk.ValAddress) ([]stakingtypes.Delegation, error)
	}

	SlashingKeeper interface {
		IsTombstoned(context.Context, sdk.ConsAddress) bool
		Tombstone(context.Context, sdk.ConsAddress) error
		Jail(context.Context, sdk.ConsAddress) error
		JailUntil(context.Context, sdk.ConsAddress, time.Time) error
	}
)
