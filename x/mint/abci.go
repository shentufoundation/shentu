package mint

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/shentufoundation/shentu/v2/x/mint/keeper"
)

// Minting module event types
const (
	EventTypeMint = "mint"

	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx context.Context, k keeper.Keeper, ic types.InflationCalculationFn) error {
	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// recalculate inflation rate
	totalBondedTokens, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	minter.Inflation = ic(ctx, minter, params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalBondedTokens)
	if err = k.Minter.Set(ctx, minter); err != nil {
		return err
	}

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	if err = k.MintCoins(ctx, mintedCoins); err != nil {
		return err
	}

	communityPoolRatio, err := k.GetCommunityPoolRatio(ctx)
	if err != nil {
		return err
	}
	communityPoolCoin, err := k.GetPoolMint(ctx, communityPoolRatio, mintedCoin)
	if err != nil {
		return err
	}
	communityPoolCoins := sdk.NewCoins(communityPoolCoin)

	collectedFeesCoins := mintedCoins.Sub(communityPoolCoins...)

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, collectedFeesCoins); err != nil {
		panic(err)
	}

	if err = k.SendToCommunityPool(ctx, communityPoolCoins); err != nil {
		panic(err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, coin := range mintedCoins {
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeMint,
				sdk.NewAttribute(AttributeKeyBondedRatio, bondedRatio.String()),
				sdk.NewAttribute(AttributeKeyInflation, minter.Inflation.String()),
				sdk.NewAttribute(AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
			),
		)
	}

	return nil
}
