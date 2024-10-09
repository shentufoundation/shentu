// Package types only consists of a copy from the cosmos' mint module's expected keepers
package types

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	SetModuleAccount(ctx context.Context, i sdk.ModuleAccountI)
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	StakingTokenSupply(context.Context) (math.Int, error) // total staking token supply
	TotalBondedTokens(context.Context) (math.Int, error)  // total bonded tokens within the validator set
	BondedRatio(ctx context.Context) (math.LegacyDec, error)
	BondDenom(ctx context.Context) (string, error)
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx context.Context, name string, amt sdk.Coins) error
}

// DistributionKeeper defines the expected distribution keeper.
type DistributionKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
