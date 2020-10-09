package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPremiumRate(days uint64) sdk.Dec {
	return sdk.NewDecFromBigIntWithPrec(big.NewInt(4), 2) //placeholder 4% for now
}

// GetLockAmount returns the proportional collateral amount
// to lock given some loss amount.
func GetLockAmount(loss, totalCollateral, collateral sdk.Int) sdk.Int {
	lossAmount := loss
	totalCollateralAmount := totalCollateral
	if totalCollateralAmount.IsZero() {
		return sdk.ZeroInt()
	}
	collateralAmount := collateral
	lockAmount := lossAmount.Mul(collateralAmount).Quo(totalCollateralAmount)
	return lockAmount
}
