package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minting module event types
const (
	EventTypeMint = "mint"

	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	totalBondedTokens := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalBondedTokens)
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)
	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	communityPoolRatio := k.GetCommunityPoolRatio(ctx)
	communityPoolCoins := k.GetCommunityPoolMint(ctx, communityPoolRatio, mintedCoin)
	collectedFeesCoins := mintedCoins.Sub(communityPoolCoins)

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, collectedFeesCoins); err != nil {
		panic(err)
	}

	if err = k.SendToCommunityPool(ctx, communityPoolCoins); err != nil {
		panic(err)
	}

	mintedCoins = mintedCoins.Add(communityPoolCoins...)
	for _, coin := range mintedCoins {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeMint,
				sdk.NewAttribute(AttributeKeyBondedRatio, bondedRatio.String()),
				sdk.NewAttribute(AttributeKeyInflation, minter.Inflation.String()),
				sdk.NewAttribute(AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
			),
		)
	}
}
