package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MicroCTKDenom = "uctk"
	MicroUnit     = int64(1e6)

	Update1Height = 348000
)

// GetCoinPercentage returns a certain percentage of coins.
// NOTE: The amount of coins returned will always be floored.
func GetCoinPercentage(coins sdk.Coins, percentage int64) sdk.Coins {
	if percentage > 100 {
		percentage = 100
	} else if percentage < 0 {
		percentage = 0
	}
	res := sdk.Coins{}
	for _, coin := range coins {
		res = res.Add(sdk.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Mul(sdk.NewInt(percentage)).Quo(sdk.NewInt(100)),
		})
	}
	return res
}

// DivideCoins divides the coins with certain number, discarding any remainders.
func DivideCoins(coins sdk.Coins, dividend int64) sdk.Coins {
	res := sdk.Coins{}
	for _, coin := range coins {
		res = res.Add(sdk.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.Quo(sdk.NewInt(dividend)),
		})
	}
	return res
}
