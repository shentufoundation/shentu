package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/common"
)

func GetPremiumRate(days uint64) sdk.Dec {
	return sdk.NewDecFromBigIntWithPrec(big.NewInt(4), 2) //placeholder 4% for now
}

func GetLockedCoins(loss, totalCollateral, collateral sdk.Coins) sdk.Coins {
	lossAmount := loss.AmountOf(common.MicroCTKDenom)
	totalCollateralAmount := totalCollateral.AmountOf(common.MicroCTKDenom)
	collateralAmount := collateral.AmountOf(common.MicroCTKDenom)
	lockedAmount := lossAmount.Mul(collateralAmount).Quo(totalCollateralAmount)
	return sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, lockedAmount))
}
