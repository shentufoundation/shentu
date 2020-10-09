package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
