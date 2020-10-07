package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPremiumRate(days uint64) sdk.Dec {
	return sdk.NewDecFromBigIntWithPrec(big.NewInt(4), 2) //placeholder 4% for now
}

func GetLockedCoins(loss, totalCollateral, collateral sdk.Coins, denom string) sdk.Coins {
	lossAmount := loss.AmountOf(denom)
	totalCollateralAmount := totalCollateral.AmountOf(denom)
	if totalCollateralAmount.IsZero() {
		return sdk.Coins{}
	}
	collateralAmount := collateral.AmountOf(denom)
	lockedAmount := lossAmount.Mul(collateralAmount).Quo(totalCollateralAmount)
	return sdk.NewCoins(sdk.NewCoin(denom, lockedAmount))
}
