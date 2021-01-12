// Package types only consists of a copy from the cosmos' mint module's expected keepers
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	StakingTokenSupply(ctx sdk.Context) sdk.Int
	TotalBondedTokens(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
	BondDenom(ctx sdk.Context) string
}

// SupplyKeeper defines the expected supply keeper.
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, exported.ModuleAccountI)
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

// DistributionKeeper defines the expected distribution keeper.
type DistributionKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
	GetFeePool(ctx sdk.Context) distribution.FeePool
}

type ShieldKeeper interface {
	GetGlobalShieldStakingPool(ctx sdk.Context) sdk.Int
	FundShieldBlockRewards(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}
